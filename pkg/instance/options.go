package instance

import (
	"encoding/json"
	"fmt"
	"llamadeck/pkg/backends"
	"llamadeck/pkg/config"
	"llamadeck/pkg/validation"
	"log"
	"maps"
	"slices"
	"sync"
)

// Options contains the actual configuration (exported - this is the public API).
type Options struct {
	// Auto restart
	AutoRestart  *bool `json:"auto_restart,omitempty"`
	MaxRestarts  *int  `json:"max_restarts,omitempty"`
	RestartDelay *int  `json:"restart_delay,omitempty"` // seconds
	// On demand start
	OnDemandStart *bool `json:"on_demand_start,omitempty"`
	// Idle timeout
	IdleTimeout *int `json:"idle_timeout,omitempty"` // minutes
	// Environment variables
	Environment map[string]string `json:"environment,omitempty"`
	// Preset configuration
	PresetIni *string `json:"preset_ini,omitempty"`

	// Execution context overrides
	DockerEnabled   *bool  `json:"docker_enabled,omitempty"`
	CommandOverride string `json:"command_override,omitempty"`

	// Instance group for hierarchical eviction
	Group string `json:"group,omitempty"`

	// Assigned nodes
	Nodes map[string]struct{} `json:"-"`
	// Backend options
	BackendOptions backends.Options `json:"-"`
}

// options wraps Options with thread-safe access (unexported).
type options struct {
	mu   sync.RWMutex
	opts *Options
}

// newOptions creates a new options wrapper with the given Options
func newOptions(opts *Options) *options {
	return &options{
		opts: opts,
	}
}

// get returns a copy of the current options
func (o *options) get() *Options {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.opts
}

// set updates the options
func (o *options) set(opts *Options) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.opts = opts
}

func (o *options) GetHost() string {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.opts.BackendOptions.GetHost()
}

func (o *options) GetPort() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.opts.BackendOptions.GetPort()
