package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func getDefaultConfig(dataDir string) AppConfig {
	return AppConfig{
		Server: ServerConfig{
			Host:           "0.0.0.0",
			Port:           8080,
			AllowedOrigins: []string{"*"}, // Default to allow all origins
			AllowedHeaders: []string{"*"}, // Default to allow all headers
			EnableSwagger:  false,
		},
		LocalNode: "main",
		Nodes:     map[string]NodeConfig{},
		DataDir:   dataDir,
		Backends: BackendConfig{
			LlamaCpp: BackendSettings{
				Command:         "llama-server",
				Args:            []string{},
				Environment:     map[string]string{},
				HealthPath:      "/health",
				CacheDir:        getDefaultLlamaCacheDir(),
				DownloadTimeout: 3600 * time.Second,
				Docker: &DockerSettings{
					Enabled: false,
					Image:   "ghcr.io/ggml-org/llama.cpp:server",
					Args: []string{
						"run", "--rm", "--network", "host", "--gpus", "all",
						"-v", getDefaultLlamaCacheDir() + ":/root/.cache/llama.cpp"},
					Environment: map[string]string{},
				},
			},
			VLLM: BackendSettings{
				Command:    "vllm",
				Args:       []string{"serve"},
				HealthPath: "/health",
				Docker: &DockerSettings{
					Enabled: false,
					Image:   "vllm/vllm-openai:latest",
					Args: []string{
						"run", "--rm", "--network", "host", "--gpus", "all", "--shm-size", "1g",
						"-v", filepath.Join(dataDir, "huggingface") + ":/root/.cache/huggingface",
					},
					Environment: map[string]string{},
				},
			},
			MLX: BackendSettings{
				Command:    "mlx_lm.server",
				Args:       []string{},
				HealthPath: "/health",
				// No Docker section for MLX - not supported
			},
			Custom: map[string]BackendSettings{},
		},
		Instances: InstancesConfig{
			PortRange:            [2]int{8000, 9000},
			AutoCreateDirs:       true,
			MaxInstances:         -1, // -1 means unlimited
			MaxRunningInstances:  -1, // -1 means unlimited
			GroupLimits:          map[string]int{},
			EnableLRUEviction:    true,
			DefaultIdleTimeout:   30, // Default idle timeout of 30 minutes
			DefaultAutoRestart:   true,
			DefaultMaxRestarts:   3,
			DefaultRestartDelay:  5,
			DefaultOnDemandStart: true,
			OnDemandStartTimeout: 120, // 2 minutes
			TimeoutCheckInterval: 5,   // Check timeouts every 5 minutes
			LogsDir:              "",  // Will be set to data_dir/logs if empty
			InstancesDir:         "",  // Will be set to data_dir/instances if empty
			LogRotationEnabled:   true,
			LogRotationMaxSize:   100,
			LogRotationCompress:  false,
		},
		Database: DatabaseConfig{
			Path:               "", // Will be set to data_dir/llamadeck.db if empty
			MaxOpenConnections: 25,
			MaxIdleConnections: 5,
			ConnMaxLifetime:    5 * time.Minute,
		},
		Auth: AuthConfig{
			RequireInferenceAuth:  true,
			RequireManagementAuth: true,
			ManagementKeys:        []string{},
			OIDC: OIDCConfig{
				Scopes:       []string{"openid", "profile", "email"},
