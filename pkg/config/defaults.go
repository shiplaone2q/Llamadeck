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
