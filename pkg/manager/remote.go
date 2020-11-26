package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"llamadeck/pkg/config"
	"llamadeck/pkg/instance"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const apiBasePath = "/api/v1/instances/"

var ErrInstanceNotFound = errors.New("instance not found")

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrInstanceNotFound) {
		return true
	}
	var rnf *RemoteNotFoundError
	return errors.As(err, &rnf)
}

// RemoteNotFoundError is returned by remoteManager CRUD functions when the
// remote node reports the target instance does not exist
type RemoteNotFoundError struct {
	Name string
	Err  error
}

func (e *RemoteNotFoundError) Error() string { return e.Err.Error() }
func (e *RemoteNotFoundError) Unwrap() error { return e.Err }

func (e *RemoteNotFoundError) Is(target error) bool {
	return target == ErrInstanceNotFound
}

// remoteManager handles HTTP operations for remote instances.
type remoteManager struct {
	mu             sync.RWMutex
	client         *http.Client
	nodeMap        map[string]*config.NodeConfig // node name -> node config
	instanceToNode map[string]*config.NodeConfig // instance name -> node config
}

// newRemoteManager creates a new remote manager.
func newRemoteManager(nodes map[string]config.NodeConfig, timeout time.Duration) *remoteManager {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	// Build node config map
	nodeMap := make(map[string]*config.NodeConfig)
	for name := range nodes {
		nodeCopy := nodes[name]
		nodeMap[name] = &nodeCopy
	}

	return &remoteManager{
		client: &http.Client{
			Timeout: timeout,
		},
		nodeMap:        nodeMap,
		instanceToNode: make(map[string]*config.NodeConfig),
	}
}

// GetNodeForInstance returns the node configuration for a given instance.
// Returns nil if the instance is not mapped to any node.
func (rm *remoteManager) getNodeForInstance(instanceName string) (*config.NodeConfig, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	node, exists := rm.instanceToNode[instanceName]
	return node, exists
}

// SetInstanceNode maps an instance to a specific node.
// Returns an error if the node doesn't exist.
func (rm *remoteManager) setInstanceNode(instanceName, nodeName string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	node, exists := rm.nodeMap[nodeName]
	if !exists {
		return fmt.Errorf("node %s not found", nodeName)
	}

	rm.instanceToNode[instanceName] = node
	return nil
}

// RemoveInstance removes the instance-to-node mapping.
func (rm *remoteManager) removeInstance(instanceName string) {
	rm.mu.Lock()
