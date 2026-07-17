package manager

import (
	"fmt"
	"math/bits"
	"sync"
)

// portAllocator provides efficient port allocation using a bitmap for O(1) operations.
// The bitmap approach prevents unbounded memory growth and simplifies port management.
type portAllocator struct {
	mu sync.Mutex

	// Bitmap for O(1) allocation/release
	// Each bit represents a port (1 = allocated, 0 = free)
	bitmap []uint64 // Each uint64 covers 64 ports

	// Map port to instance name for cleanup operations
	allocated map[int]string

	minPort   int
	maxPort   int
	rangeSize int
}

// newPortAllocator creates a new port allocator for the given port range.
func newPortAllocator(minPort, maxPort int) *portAllocator {
	rangeSize := maxPort - minPort + 1
	bitmapSize := (rangeSize + 63) / 64 // Round up to nearest uint64

	return &portAllocator{
		bitmap:    make([]uint64, bitmapSize),
		allocated: make(map[int]string),
		minPort:   minPort,
		maxPort:   maxPort,
		rangeSize: rangeSize,
	}
}

// allocate finds and allocates the first available port for the given instance.
// Returns the allocated port or an error if no ports are available.
func (p *portAllocator) allocate(instanceName string) (int, error) {
	if instanceName == "" {
		return 0, fmt.Errorf("instance name cannot be empty")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	port, err := p.findFirstFreeBit()
	if err != nil {
		return 0, err
	}

	p.setBit(port)
	p.allocated[port] = instanceName

	return port, nil
}

// allocateSpecific allocates a specific port for the given instance.
// Returns an error if the port is already allocated or out of range.
func (p *portAllocator) allocateSpecific(port int, instanceName string) error {
	if instanceName == "" {
		return fmt.Errorf("instance name cannot be empty")
	}
	if port < p.minPort || port > p.maxPort {
		return fmt.Errorf("port %d is out of range [%d-%d]", port, p.minPort, p.maxPort)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isBitSet(port) {
		return fmt.Errorf("port %d is already allocated", port)
	}

	p.setBit(port)
	p.allocated[port] = instanceName

	return nil
}

// release releases a specific port, making it available for reuse.
// Returns an error if the port is not allocated.
func (p *portAllocator) release(port int) error {
	if port < p.minPort || port > p.maxPort {
		return fmt.Errorf("port %d is out of range [%d-%d]", port, p.minPort, p.maxPort)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isBitSet(port) {
		return fmt.Errorf("port %d is not allocated", port)
	}

	p.clearBit(port)
	delete(p.allocated, port)

	return nil
}

// releaseByInstance releases all ports allocated to the given instance.
