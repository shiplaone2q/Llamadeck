package instance_test

import (
	"encoding/json"
	"llamadeck/pkg/backends"
	"llamadeck/pkg/config"
	"llamadeck/pkg/instance"
	"llamadeck/pkg/testutil"
	"os"
	"testing"
	"time"
)

func TestNewInstance(t *testing.T) {
	globalConfig := &config.AppConfig{
		Backends: config.BackendConfig{
			LlamaCpp: config.BackendSettings{
				Command: "llama-server",
				Args:    []string{},
			},
			MLX: config.BackendSettings{
				Command: "mlx_lm.server",
				Args:    []string{},
			},
			VLLM: config.BackendSettings{
				Command: "vllm",
				Args:    []string{"serve"},
			},
		},
		Instances: config.InstancesConfig{
			DefaultAutoRestart:  true,
			LogsDir:             "/tmp/test",
			DefaultMaxRestarts:  3,
			DefaultRestartDelay: 5,
			DefaultIdleTimeout:  30,
		},
		Nodes:     map[string]config.NodeConfig{},
		LocalNode: "main",
	}

	options := &instance.Options{
		BackendOptions: backends.Options{
			BackendType: backends.BackendTypeLlamaCpp,
			LlamaServerOptions: &backends.LlamaServerOptions{
				Model: "/path/to/model.gguf",
				Port:  8080,
			},
		},
	}

	// Mock onStatusChange function
	mockOnStatusChange := func(oldStatus, newStatus instance.Status) {}

	inst := instance.New("test-instance", globalConfig, options, mockOnStatusChange)

	if inst.Name != "test-instance" {
		t.Errorf("Expected name 'test-instance', got %q", inst.Name)
	}
	if inst.IsRunning() {
		t.Error("New instance should not be running")
	}

	// Check that options were properly set with defaults applied
	opts := inst.GetOptions()
	if opts.BackendOptions.LlamaServerOptions.Model != "/path/to/model.gguf" {
		t.Errorf("Expected model '/path/to/model.gguf', got %q", opts.BackendOptions.LlamaServerOptions.Model)
	}
	if inst.GetPort() != 8080 {
		t.Errorf("Expected port 8080, got %d", inst.GetPort())
	}

	// Check that defaults were applied
	if opts.AutoRestart == nil || !*opts.AutoRestart {
		t.Error("Expected AutoRestart to be true (default)")
	}
	if opts.MaxRestarts == nil || *opts.MaxRestarts != 3 {
		t.Errorf("Expected MaxRestarts to be 3 (default), got %v", opts.MaxRestarts)
	}
	if opts.RestartDelay == nil || *opts.RestartDelay != 5 {
		t.Errorf("Expected RestartDelay to be 5 (default), got %v", opts.RestartDelay)
	}
	if opts.IdleTimeout == nil || *opts.IdleTimeout != 30 {
		t.Errorf("Expected IdleTimeout to be 30 (default), got %v", opts.IdleTimeout)
	}

	// Test that explicit values override defaults
	autoRestart := false
	maxRestarts := 10
	optionsWithOverrides := &instance.Options{
		AutoRestart: &autoRestart,
		MaxRestarts: &maxRestarts,
		BackendOptions: backends.Options{
			BackendType: backends.BackendTypeLlamaCpp,
			LlamaServerOptions: &backends.LlamaServerOptions{
				Model: "/path/to/model.gguf",
			},
		},
	}

	inst2 := instance.New("test-override", globalConfig, optionsWithOverrides, mockOnStatusChange)
	opts2 := inst2.GetOptions()

	if opts2.AutoRestart == nil || *opts2.AutoRestart {
		t.Error("Expected AutoRestart to be false (overridden)")
	}
	if opts2.MaxRestarts == nil || *opts2.MaxRestarts != 10 {
		t.Errorf("Expected MaxRestarts to be 10 (overridden), got %v", opts2.MaxRestarts)
	}
}

func TestSetOptions(t *testing.T) {
	globalConfig := &config.AppConfig{
		Backends: config.BackendConfig{
			LlamaCpp: config.BackendSettings{
				Command: "llama-server",
				Args:    []string{},
			},
			MLX: config.BackendSettings{
				Command: "mlx_lm.server",
				Args:    []string{},
			},
			VLLM: config.BackendSettings{
				Command: "vllm",
				Args:    []string{"serve"},
			},
		},
		Instances: config.InstancesConfig{
			DefaultAutoRestart:  true,
			LogsDir:             "/tmp/test",
			DefaultMaxRestarts:  3,
			DefaultRestartDelay: 5,
			DefaultIdleTimeout:  30,
		},
		Nodes:     map[string]config.NodeConfig{},
		LocalNode: "main",
	}

	initialOptions := &instance.Options{
		BackendOptions: backends.Options{
			BackendType: backends.BackendTypeLlamaCpp,
			LlamaServerOptions: &backends.LlamaServerOptions{
				Model: "/path/to/model.gguf",
				Port:  8080,
			},
		},
	}

	// Mock onStatusChange function
	mockOnStatusChange := func(oldStatus, newStatus instance.Status) {}

	inst := instance.New("test-instance", globalConfig, initialOptions, mockOnStatusChange)

	// Update options
	newOptions := &instance.Options{
		BackendOptions: backends.Options{
			BackendType: backends.BackendTypeLlamaCpp,
			LlamaServerOptions: &backends.LlamaServerOptions{
				Model: "/path/to/new-model.gguf",
				Port:  8081,
			},
		},
	}

	inst.SetOptions(newOptions)
	opts := inst.GetOptions()

	if opts.BackendOptions.LlamaServerOptions.Model != "/path/to/new-model.gguf" {
		t.Errorf("Expected updated model '/path/to/new-model.gguf', got %q", opts.BackendOptions.LlamaServerOptions.Model)
	}
	if inst.GetPort() != 8081 {
		t.Errorf("Expected updated port 8081, got %d", inst.GetPort())
	}

	// Check that defaults are still applied
	if opts.AutoRestart == nil || !*opts.AutoRestart {
		t.Error("Expected AutoRestart to be true (default)")
	}
}

func TestMarshalJSON(t *testing.T) {
	globalConfig := &config.AppConfig{
		Backends: config.BackendConfig{
			LlamaCpp: config.BackendSettings{Command: "llama-server"},
		},
		Instances: config.InstancesConfig{LogsDir: "/tmp/test"},
		Nodes:     map[string]config.NodeConfig{},
		LocalNode: "main",
	}
	options := &instance.Options{
		BackendOptions: backends.Options{
			BackendType: backends.BackendTypeLlamaCpp,
			LlamaServerOptions: &backends.LlamaServerOptions{
				Model: "/path/to/model.gguf",
				Port:  8080,
			},
		},
	}

	inst := instance.New("test-instance", globalConfig, options, nil)

	data, err := json.Marshal(inst)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Verify by unmarshaling and checking key fields
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if result["name"] != "test-instance" {
		t.Errorf("Expected name 'test-instance', got %v", result["name"])
	}
	if result["status"] != "stopped" {
		t.Errorf("Expected status 'stopped', got %v", result["status"])
	}
	if result["options"] == nil {
		t.Error("Expected options to be included in JSON")
	}
}

func TestUnmarshalJSON(t *testing.T) {
	jsonData := `{
		"name": "test-instance",
		"status": "running",
		"options": {
			"auto_restart": false,
			"max_restarts": 5,
			"backend_type": "llama_cpp",
			"backend_options": {
				"model": "/path/to/model.gguf",
				"port": 8080
