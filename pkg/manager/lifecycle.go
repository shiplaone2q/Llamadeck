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

	return &lifecycleManager{
		registry:      registry,
		manager:       manager,
		ticker:        time.NewTicker(checkInterval),
		checkInterval: checkInterval,
		enableLRU:     enableLRU,
		shutdownChan:  make(chan struct{}),
		shutdownDone:  make(chan struct{}),
	}
}

// Start begins the timeout checking loop in a goroutine.
func (l *lifecycleManager) start() {
	go l.timeoutCheckLoop()
}

// Stop gracefully stops the lifecycle manager.
// This ensures the timeout checker completes before instance cleanup begins.
func (l *lifecycleManager) stop() {
	l.shutdownOnce.Do(func() {
		close(l.shutdownChan)
		<-l.shutdownDone // Wait for checker to finish (prevents shutdown race)
		l.ticker.Stop()
	})
}

// timeoutCheckLoop is the main loop that periodically checks for timeouts.
func (l *lifecycleManager) timeoutCheckLoop() {
	defer close(l.shutdownDone) // Signal completion

	for {
		select {
		case <-l.ticker.C:
			l.checkTimeouts()
		case <-l.shutdownChan:
