package manager_test

import (
	"llamadeck/pkg/backends"
	"llamadeck/pkg/instance"
	"llamadeck/pkg/manager"
	"sync"
	"testing"
	"time"
)

func TestInstanceTimeoutLogic(t *testing.T) {
	testManager := createTestManager(t)
	defer testManager.Shutdown()

	idleTimeout := 1 // 1 minute
	inst := createInstanceWithTimeout(t, testManager, "timeout-test", "/path/to/model.gguf", &idleTimeout)

	// Test timeout logic with mock time provider
	mockTime := NewMockTimeProvider(time.Now())
	inst.SetTimeProvider(mockTime)

	// Set instance to running state so timeout logic can work
	inst.SetStatus(instance.Running)
	defer inst.SetStatus(instance.Stopped)

	// Update last request time
	inst.UpdateLastRequestTime()

	// Initially should not timeout (just updated)
	if inst.ShouldTimeout() {
		t.Error("Instance should not timeout immediately after request")
	}

	// Advance time to trigger timeout
	mockTime.SetTime(time.Now().Add(2 * time.Minute))

	// Now it should timeout
	if !inst.ShouldTimeout() {
		t.Error("Instance should timeout after idle period")
	}
}

func TestInstanceWithoutTimeoutNeverExpires(t *testing.T) {
	testManager := createTestManager(t)
	defer testManager.Shutdown()

	zeroTimeout := 0
	noTimeoutInst := createInstanceWithTimeout(t, testManager, "no-timeout-test", "/path/to/model.gguf", &zeroTimeout)

	mockTime := NewMockTimeProvider(time.Now())
	noTimeoutInst.SetTimeProvider(mockTime)
	noTimeoutInst.SetStatus(instance.Running)
	defer noTimeoutInst.SetStatus(instance.Stopped)

	noTimeoutInst.UpdateLastRequestTime()

	// Advance time significantly
	mockTime.SetTime(mockTime.Now().Add(24 * time.Hour))

	// Even with time advanced, should not timeout
	if noTimeoutInst.ShouldTimeout() {
		t.Error("Instance without timeout configuration should never timeout")
	}
}

func TestEvictLRUInstance_Success(t *testing.T) {
	manager := createTestManager(t)
	defer manager.Shutdown()

	// Create 3 instances with idle timeout enabled (value doesn't matter for LRU logic)
	validTimeout := 1
	inst1 := createInstanceWithTimeout(t, manager, "instance-1", "/path/to/model1.gguf", &validTimeout)
	inst2 := createInstanceWithTimeout(t, manager, "instance-2", "/path/to/model2.gguf", &validTimeout)
	inst3 := createInstanceWithTimeout(t, manager, "instance-3", "/path/to/model3.gguf", &validTimeout)

	// Set up mock time and set instances to running
	mockTime := NewMockTimeProvider(time.Now())
	inst1.SetTimeProvider(mockTime)
	inst2.SetTimeProvider(mockTime)
	inst3.SetTimeProvider(mockTime)

	inst1.SetStatus(instance.Running)
	inst2.SetStatus(instance.Running)
	inst3.SetStatus(instance.Running)
	defer func() {
		// Clean up - ensure all instances are stopped
