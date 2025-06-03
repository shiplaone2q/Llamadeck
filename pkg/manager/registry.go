package manager

import (
	"fmt"
	"llamadeck/pkg/instance"
	"sync"
)

// instanceRegistry provides thread-safe storage and lookup of instances
// with running state tracking using lock-free sync.Map for status checks.
type instanceRegistry struct {
	mu        sync.RWMutex
	instances map[string]*instance.Instance
	running   sync.Map // map[string]struct{} - lock-free for status checks
}

// newInstanceRegistry creates a new instance registry.
func newInstanceRegistry() *instanceRegistry {
	return &instanceRegistry{
		instances: make(map[string]*instance.Instance),
	}
}

// Get retrieves an instance by name.
// Returns the instance and true if found, nil and false otherwise.
func (r *instanceRegistry) get(name string) (*instance.Instance, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	inst, exists := r.instances[name]
	return inst, exists
}

// List returns a snapshot copy of all instances to prevent external mutation.
func (r *instanceRegistry) list() []*instance.Instance {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*instance.Instance, 0, len(r.instances))
	for _, inst := range r.instances {
		result = append(result, inst)
	}
	return result
}

// ListRunning returns a snapshot of all currently running instances.
func (r *instanceRegistry) listRunning() []*instance.Instance {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*instance.Instance, 0)
	for name, inst := range r.instances {
		if _, isRunning := r.running.Load(name); isRunning {
			result = append(result, inst)
		}
	}
	return result
}

// Add adds a new instance to the registry.
// Returns an error if an instance with the same name already exists.
func (r *instanceRegistry) add(inst *instance.Instance) error {
	if inst == nil {
		return fmt.Errorf("cannot add nil instance")
	}

