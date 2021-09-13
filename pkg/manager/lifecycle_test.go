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
