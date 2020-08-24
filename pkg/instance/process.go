package instance

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

// process manages the OS process lifecycle for a local instance.
// process owns its complete lifecycle including auto-restart logic.
type process struct {
	instance *Instance // Back-reference for SetStatus, GetOptions

	mu            sync.RWMutex
	cmd           *exec.Cmd
	ctx           context.Context
	cancel        context.CancelFunc
	stdin         io.Closer
	stdout        io.ReadCloser
	stderr        io.ReadCloser
	restarts      int
	restartCancel context.CancelFunc
	monitorDone   chan struct{}
}

// newProcess creates a new process component for the given instance
func newProcess(instance *Instance) *process {
	return &process{
		instance: instance,
	}
}

// start starts the OS process and returns an error if it fails.
func (p *process) start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.instance.IsRunning() {
		return fmt.Errorf("instance %s is already running", p.instance.Name)
	}

	// Safety check: ensure options are valid
	if p.instance.options == nil {
		return fmt.Errorf("instance %s has no options set", p.instance.Name)
	}

	// Reset restart counter when manually starting (not during auto-restart)
	// We can detect auto-restart by checking if restartCancel is set
	if p.restartCancel == nil {
		p.restarts = 0
	}

	// Initialize last request time to current time when starting
	if p.instance.proxy != nil {
		p.instance.proxy.updateLastRequestTime()
	}

	// Create context before building command (needed for CommandContext)
	p.ctx, p.cancel = context.WithCancel(context.Background())

	// Create log files
	if err := p.instance.logger.create(); err != nil {
		return fmt.Errorf("failed to create log files: %w", err)
	}

	// Build command using backend-specific methods
	cmd, cmdErr := p.buildCommand()
	if cmdErr != nil {
		return fmt.Errorf("failed to build command: %w", cmdErr)
	}
	p.cmd = cmd

	// Double-spawn guard: if our port is already accepting connections, a
	// surviving llama-server from a previous (or crashed) start is still
	// bound to it. Waiting it out avoids double-binding, which would run
	// two copies of the model in VRAM at once.
	host, port := p.instance.options.GetHost(), p.instance.options.GetPort()
	// Dial to loopback regardless of the configured bind host: a server bound
	// to 0.0.0.0/:: is reachable via 127.0.0.1, but net.DialTimeout to
	// 0.0.0.0 (or ::) fails on Windows with an invalid-address error and the
	// guard would silently no-op.
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}
	if port > 0 {
		deadline := time.Now().Add(15 * time.Second)
		for portInUse(host, port) {
			if time.Now().After(deadline) {
				p.instance.logger.close()
				if p.cancel != nil {
					p.cancel()
				}
				return fmt.Errorf("port %d is still occupied by a surviving process; refusing to start a second server on it", port)
			}
			log.Printf("Instance %s: port %d still in use, waiting for it to be released...", p.instance.Name, port)
			time.Sleep(500 * time.Millisecond)
		}
	}

	setProcAttrs(p.cmd)

