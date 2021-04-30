package instance

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"llamadeck/pkg/backends"
	"llamadeck/pkg/config"
)

// Instance represents a running instance of llama server
type Instance struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Created int64  `json:"created,omitempty"` // Unix timestamp when instance was created

	// Global configuration
	globalInstanceSettings *config.InstancesConfig
	globalBackendSettings  *config.BackendConfig
	globalNodesConfig      map[string]config.NodeConfig
	localNodeName          string `json:"-"` // Name of the local node for remote detection

	status  *status  `json:"-"`
	options *options `json:"-"`

	// Components (can be nil for remote instances)
	process *process `json:"-"`
	proxy   *proxy   `json:"-"`
	logger  *logger  `json:"-"`
}

// New creates a new instance with the given name, log path, options and local node name
func New(name string, globalConfig *config.AppConfig, opts *Options, onStatusChange func(oldStatus, newStatus Status)) *Instance {

	globalInstanceSettings := &globalConfig.Instances
	globalBackendSettings := &globalConfig.Backends
	globalNodesConfig := globalConfig.Nodes
	localNodeName := globalConfig.LocalNode

	// Validate and copy options
	opts.validateAndApplyDefaults(name, globalInstanceSettings)

	// Create status wrapper
	status := newStatus(Stopped)
	status.onStatusChange = onStatusChange

	// Create options wrapper
	options := newOptions(opts)

	instance := &Instance{
		ID:                     0, // Will be set by database
		Name:                   name,
		options:                options,
		globalInstanceSettings: globalInstanceSettings,
		globalBackendSettings:  globalBackendSettings,
		globalNodesConfig:      globalNodesConfig,
		localNodeName:          localNodeName,
		Created:                time.Now().Unix(),
		status:                 status,
	}

	var err error
	instance.proxy, err = newProxy(instance)
	if err != nil {
		log.Println("Warning: Failed to create proxy for instance", instance.Name, "-", err)
	}

	// Only create logger, proxy, and process for local instances
	if !instance.IsRemote() {
		logRotationConfig := &LogRotationConfig{
			Enabled:  globalInstanceSettings.LogRotationEnabled,
			MaxSize:  globalInstanceSettings.LogRotationMaxSize,
			Compress: globalInstanceSettings.LogRotationCompress,
		}
		instance.logger = newLogger(
			name,
			globalInstanceSettings.LogsDir,
			logRotationConfig,
		)
		instance.process = newProcess(instance)

		if err := syncPresetIni(name, opts, globalInstanceSettings.InstancesDir); err != nil {
			log.Printf("Warning: Failed to write preset.ini for instance %s: %v", name, err)
		}
	}

	return instance
}

// syncPresetIni keeps the preset.ini file in sync with options:
// written when a preset is set, removed when it is cleared
func syncPresetIni(name string, opts *Options, instancesDir string) error {
	if opts == nil || opts.PresetIni == nil || *opts.PresetIni == "" {
		// A stale file would still be handed to the backend via --models-preset
		presetPath := filepath.Join(instancesDir, name, "preset.ini")
		if err := os.Remove(presetPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove preset.ini: %w", err)
		}
		return nil
	}

	instanceDir := filepath.Join(instancesDir, name)
	if err := os.MkdirAll(instanceDir, 0755); err != nil {
		return fmt.Errorf("failed to create instance directory: %w", err)
	}

	presetPath := filepath.Join(instanceDir, "preset.ini")
	if err := os.WriteFile(presetPath, []byte(*opts.PresetIni), 0644); err != nil {
		return fmt.Errorf("failed to write preset.ini: %w", err)
	}

	log.Printf("Wrote preset.ini for instance %s at %s", name, presetPath)
	return nil
}
