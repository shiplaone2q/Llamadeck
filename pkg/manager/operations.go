package manager

import (
	"context"
	"fmt"
	"llamadeck/pkg/backends"
	"llamadeck/pkg/instance"
	"llamadeck/pkg/validation"
	"log"
	"slices"
	"strings"
)

type MaxRunningInstancesError error

// updateLocalInstanceFromRemote updates the local stub instance with data from the remote instance
func (im *instanceManager) updateLocalInstanceFromRemote(localInst *instance.Instance, remoteInst *instance.Instance) {
	if localInst == nil || remoteInst == nil {
		return
	}

	remoteOptions := remoteInst.GetOptions()
	if remoteOptions == nil {
		return
	}

	// Update the local instance with all remote data
	localInst.SetOptions(remoteOptions)
	localInst.SetStatus(remoteInst.GetStatus())
	localInst.Created = remoteInst.Created
}

// ListInstances returns a list of all instances managed by the instance manager.
// For remote instances, this fetches the live state from remote nodes and updates local stubs.
func (im *instanceManager) ListInstances() ([]*instance.Instance, error) {
	instances := im.registry.list()

	// Update remote instances with live state
	ctx := context.Background()
	for _, inst := range instances {
		if node := im.getNodeForInstance(inst); node != nil {
			remoteInst, err := im.remote.getInstance(ctx, node, inst.Name)
			if err != nil {
				// Don't fail the entire list operation due to one remote failure
				continue
			}

			// Update the local stub with all remote data (preserving Nodes)
			im.updateLocalInstanceFromRemote(inst, remoteInst)
		}
	}

	return instances, nil
}

// CreateInstance creates a new instance with the given options and returns it.
// The instance is initially in a "stopped" state.
func (im *instanceManager) CreateInstance(name string, options *instance.Options) (*instance.Instance, error) {
	if options == nil {
		return nil, fmt.Errorf("instance options cannot be nil")
	}

	err := options.BackendOptions.ValidateInstanceOptions()
	if err != nil {
		return nil, err
	}

	if err := im.validateCustomBackend(options); err != nil {
		return nil, err
	}

	// Check if instance with this name already exists (must be globally unique)
	if _, exists := im.registry.get(name); exists {
		return nil, fmt.Errorf("instance with name %s already exists", name)
	}

	// Check if this is a remote instance (local node not in the Nodes set)
	if _, isLocal := options.Nodes[im.globalConfig.LocalNode]; !isLocal && len(options.Nodes) > 0 {
		// Get the first node from the set
		var nodeName string
		for node := range options.Nodes {
			nodeName = node
			break
		}

		// Create the remote instance on the remote node
		ctx := context.Background()
		nodeConfig, exists := im.remote.getNodeForInstance(nodeName)
		if !exists {
			// Try to set the node if it doesn't exist yet
			if err := im.remote.setInstanceNode(name, nodeName); err != nil {
				return nil, fmt.Errorf("node %s not found", nodeName)
			}
			nodeConfig, _ = im.remote.getNodeForInstance(name)
		}

		remoteInst, err := im.remote.createInstance(ctx, nodeConfig, name, options)
		if err != nil {
			return nil, err
		}

		// Create a local stub that preserves the Nodes field for tracking
		// We keep the original options (with Nodes) so IsRemote() works correctly
		inst := instance.New(name, im.globalConfig, options, nil)

		// Update the local stub with all remote data (preserving Nodes)
		im.updateLocalInstanceFromRemote(inst, remoteInst)

		// Map instance to node
		if err := im.remote.setInstanceNode(name, nodeName); err != nil {
			return nil, fmt.Errorf("failed to map instance to node: %w", err)
		}

		// Add to registry (doesn't count towards local limits)
		if err := im.registry.add(inst); err != nil {
			return nil, fmt.Errorf("failed to add instance to registry: %w", err)
		}

		// Persist the remote instance locally for tracking across restarts
		if err := im.persistInstance(inst); err != nil {
			// Rollback: remove from registry
			im.registry.remove(name)
			return nil, fmt.Errorf("failed to persist remote instance %s: %w", name, err)
		}

		return inst, nil
	}

	// Local instance creation
	// Check max instances limit for local instances only
	totalInstances := im.registry.count()
	remoteCount := 0
	for _, inst := range im.registry.list() {
		if inst.IsRemote() {
			remoteCount++
		}
	}
	localInstanceCount := totalInstances - remoteCount
	if localInstanceCount >= im.globalConfig.Instances.MaxInstances && im.globalConfig.Instances.MaxInstances != -1 {
		return nil, fmt.Errorf("maximum number of instances (%d) reached", im.globalConfig.Instances.MaxInstances)
	}

	// Assign and validate port for backend-specific options
	currentPort := im.getPortFromOptions(options)
	var allocatedPort int
	if currentPort == 0 {
		// Allocate a port if not specified
		allocatedPort, err = im.ports.allocate(name)
		if err != nil {
			return nil, fmt.Errorf("failed to allocate port: %w", err)
		}
		im.setPortInOptions(options, allocatedPort)
	} else {
		// Use the specified port
		if err := im.ports.allocateSpecific(currentPort, name); err != nil {
			return nil, fmt.Errorf("port %d is already in use: %w", currentPort, err)
		}
		allocatedPort = currentPort
	}

	statusCallback := func(oldStatus, newStatus instance.Status) {
		im.onStatusChange(name, oldStatus, newStatus)
	}

	inst := instance.New(name, im.globalConfig, options, statusCallback)

	// Add to registry
	if err := im.registry.add(inst); err != nil {
		// Rollback: release port
		im.ports.release(allocatedPort)
		return nil, fmt.Errorf("failed to add instance to registry: %w", err)
	}

	// Persist instance (best-effort, don't fail if persistence fails)
	if err := im.persistInstance(inst); err != nil {
		log.Printf("Warning: failed to persist instance %s: %v", name, err)
	}

	return inst, nil
}

// GetInstance retrieves an instance by its name.
// For remote instances, this fetches the live state from the remote node and updates the local stub.
//
// If the local registry has the instance but the remote reports it as gone,
// the local tracking (registry, persistence, node map) is cleaned up and
// ErrInstanceNotFound is returned. If the remote fails for any other reason
// the cached local stub is returned so callers keep working through brief
// remote outages.
func (im *instanceManager) GetInstance(name string) (*instance.Instance, error) {
	inst, exists := im.registry.get(name)
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrInstanceNotFound, name)
	}

