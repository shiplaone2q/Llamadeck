package config_test

import (
	"llamadeck/pkg/config"
	"os"
	"path/filepath"
	"testing"
)

// GetBackendSettings resolves backend settings
func getBackendSettings(bc *config.BackendConfig, backendType string) config.BackendSettings {
	switch backendType {
	case "llama-cpp":
		return bc.LlamaCpp
	case "vllm":
		return bc.VLLM
	case "mlx":
		return bc.MLX
	default:
		return config.BackendSettings{}
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Test loading config when no file exists and no env vars set
	cfg, err := config.LoadConfig("nonexistent-file.yaml")
	if err != nil {
		t.Fatalf("LoadConfig should not error with defaults: %v", err)
	}

	// Verify default values
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected default host to be 0.0.0.0, got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected default port to be 8080, got %d", cfg.Server.Port)
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	if cfg.Instances.LogsDir != filepath.Join(homedir, ".local", "share", "llamadeck", "logs") {
		t.Errorf("Expected default logs directory '%s', got %q", filepath.Join(homedir, ".local", "share", "llamadeck", "logs"), cfg.Instances.LogsDir)
	}
	if !cfg.Instances.AutoCreateDirs {
		t.Error("Expected default instances auto-create to be true")
	}
	if cfg.Instances.PortRange != [2]int{8000, 9000} {
		t.Errorf("Expected default port range [8000, 9000], got %v", cfg.Instances.PortRange)
	}
	if cfg.Instances.MaxInstances != -1 {
		t.Errorf("Expected default max instances -1, got %d", cfg.Instances.MaxInstances)
	}
	if !cfg.Instances.DefaultAutoRestart {
		t.Error("Expected default auto restart to be true")
	}
	if cfg.Instances.DefaultIdleTimeout != 30 {
		t.Errorf("Expected default idle timeout 30, got %d", cfg.Instances.DefaultIdleTimeout)
	}
	if cfg.Instances.DefaultMaxRestarts != 3 {
		t.Errorf("Expected default max restarts 3, got %d", cfg.Instances.DefaultMaxRestarts)
	}
	if cfg.Instances.DefaultRestartDelay != 5 {
		t.Errorf("Expected default restart delay 5, got %d", cfg.Instances.DefaultRestartDelay)
	}
	if cfg.Backends.LlamaCpp.HealthPath != "/health" || cfg.Backends.VLLM.HealthPath != "/health" || cfg.Backends.MLX.HealthPath != "/health" {
		t.Errorf("Expected default health path /health for all built-in backends, got llama=%q vllm=%q mlx=%q",
			cfg.Backends.LlamaCpp.HealthPath, cfg.Backends.VLLM.HealthPath, cfg.Backends.MLX.HealthPath)
	}
}

func TestLoadConfig_FromFile(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test-config.yaml")

	configContent := `
server:
  host: "localhost"
  port: 9090
instances:
  port_range: [7000, 8000]
  max_instances: 5
  logs_dir: "/custom/logs"
  llama_executable: "/usr/bin/llama-server"
  default_auto_restart: false
  default_max_restarts: 10
  default_restart_delay: 30
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify values from file
	if cfg.Server.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("Expected port 9090, got %d", cfg.Server.Port)
	}
	if cfg.Instances.PortRange != [2]int{7000, 8000} {
		t.Errorf("Expected port range [7000, 8000], got %v", cfg.Instances.PortRange)
	}
	if cfg.Instances.LogsDir != "/custom/logs" {
		t.Errorf("Expected logs directory '/custom/logs', got %q", cfg.Instances.LogsDir)
	}
	if cfg.Instances.MaxInstances != 5 {
		t.Errorf("Expected max instances 5, got %d", cfg.Instances.MaxInstances)
	}
	if cfg.Instances.DefaultAutoRestart {
		t.Error("Expected auto restart to be false")
	}
	if cfg.Instances.DefaultMaxRestarts != 10 {
		t.Errorf("Expected max restarts 10, got %d", cfg.Instances.DefaultMaxRestarts)
	}
	if cfg.Instances.DefaultRestartDelay != 30 {
		t.Errorf("Expected restart delay 30, got %d", cfg.Instances.DefaultRestartDelay)
	}
}

func TestLoadConfig_CustomBackends(t *testing.T) {
	cfg, err := config.LoadConfig("nonexistent-file.yaml")
	if err != nil {
		t.Fatalf("LoadConfig should not error with defaults: %v", err)
	}
	if cfg.Backends.Custom == nil {
		t.Fatal("Expected default Custom map to be non-nil")
	}
	if len(cfg.Backends.Custom) != 0 {
		t.Errorf("Expected default Custom map to be empty, got %v", cfg.Backends.Custom)
	}

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "custom-backends.yaml")

	configContent := `
