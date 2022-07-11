package manager

import (
	"context"
	"fmt"
	"llamadeck/pkg/config"
	"llamadeck/pkg/database"
	"llamadeck/pkg/instance"
	"log"
	"sort"
	"sync"
	"time"
)

// InstanceManager defines the interface for managing instances of the llama server.
type InstanceManager interface {
	ListInstances() ([]*instance.Instance, error)
	CreateInstance(name string, options *instance.Options) (*instance.Instance, error)
	GetInstance(name string) (*instance.Instance, error)
	UpdateInstance(name string, options *instance.Options) (*instance.Instance, error)
	DeleteInstance(name string) error
	StartInstance(name string) (*instance.Instance, error)
	AtMaxRunning() bool
	CountRunningInGroup(group string) int
	StopInstance(name string) (*instance.Instance, error)
	EvictLRUInstance(group string) error
	RestartInstance(name string) (*instance.Instance, error)
	GetInstanceLogs(name string, numLines int) (string, error)
	Shutdown()
}

type instanceManager struct {
	// Components (each with own synchronization)
	registry  *instanceRegistry
	ports     *portAllocator
	db        database.InstanceStore
	remote    *remoteManager
	lifecycle *lifecycleManager

	// Configuration
	globalConfig *config.AppConfig

	// Synchronization
	instanceLocks sync.Map // map[string]*sync.Mutex - per-instance locks for concurrent operations
	shutdownOnce  sync.Once
}

// New creates a new instance of InstanceManager with dependency injection.
func New(globalConfig *config.AppConfig, db database.InstanceStore) InstanceManager {

	if globalConfig.Instances.TimeoutCheckInterval <= 0 {
		globalConfig.Instances.TimeoutCheckInterval = 5 // Default to 5 minutes if not set
	}

	// Initialize components
	registry := newInstanceRegistry()

	// Initialize port allocator
	portRange := globalConfig.Instances.PortRange
	ports := newPortAllocator(portRange[0], portRange[1])

	// Initialize remote manager
	remote := newRemoteManager(globalConfig.Nodes, 30*time.Second)

	// Create manager instance
	im := &instanceManager{
		registry:     registry,
		ports:        ports,
		db:           db,
		remote:       remote,
		globalConfig: globalConfig,
	}

	// Initialize lifecycle manager (needs reference to manager for Stop/Evict operations)
	checkInterval := time.Duration(globalConfig.Instances.TimeoutCheckInterval) * time.Minute
	im.lifecycle = newLifecycleManager(registry, im, checkInterval, true)

	// Load existing instances from disk
	if err := im.loadInstances(); err != nil {
		log.Printf("Error loading instances: %v", err)
	}

	// Start the lifecycle manager
	im.lifecycle.start()

	return im
}

// persistInstance saves an instance using the persistence layer
func (im *instanceManager) persistInstance(inst *instance.Instance) error {
	return im.db.Save(inst)
}

func (im *instanceManager) Shutdown() {
	im.shutdownOnce.Do(func() {
		// 1. Stop lifecycle manager (stops timeout checker)
		im.lifecycle.stop()

		// 2. Get running instances (no lock needed - registry handles it)
		running := im.registry.listRunning()

		// 3. Stop local instances concurrently
		var wg sync.WaitGroup
		for _, inst := range running {
			if inst.IsRemote() {
				continue // Skip remote instances
			}
			wg.Add(1)
			go func(inst *instance.Instance) {
				defer wg.Done()
				log.Printf("Stopping instance %s...\n", inst.Name)
				if err := inst.Stop(); err != nil {
					log.Printf("Error stopping instance %s: %v\n", inst.Name, err)
				}
			}(inst)
		}
