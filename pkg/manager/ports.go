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
