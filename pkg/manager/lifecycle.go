package manager

import (
	"fmt"
	"llamadeck/pkg/instance"
	"log"
	"sync"
	"time"
)

// lifecycleManager handles background timeout checking and LRU eviction.
// It properly coordinates shutdown to prevent races with the timeout checker.
type lifecycleManager struct {
	registry *instanceRegistry
	manager  InstanceManager // For calling Stop/Evict operations

	ticker        *time.Ticker
	checkInterval time.Duration
	enableLRU     bool

	shutdownChan chan struct{}
	shutdownDone chan struct{}
	shutdownOnce sync.Once
}

// newLifecycleManager creates a new lifecycle manager.
func newLifecycleManager(
	registry *instanceRegistry,
	manager InstanceManager,
	checkInterval time.Duration,
	enableLRU bool,
) *lifecycleManager {
	if checkInterval <= 0 {
		checkInterval = 5 * time.Minute // Default to 5 minutes
	}

