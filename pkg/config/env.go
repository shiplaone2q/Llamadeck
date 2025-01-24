package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// loadEnvVars overrides config with environment variables
func loadEnvVars(cfg *AppConfig) {
	// Server config
	if host := os.Getenv("LLAMADECK_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("LLAMADECK_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}
	if allowedOrigins := os.Getenv("LLAMADECK_ALLOWED_ORIGINS"); allowedOrigins != "" {
		cfg.Server.AllowedOrigins = strings.Split(allowedOrigins, ",")
	}
	if enableSwagger := os.Getenv("LLAMADECK_ENABLE_SWAGGER"); enableSwagger != "" {
		if b, err := strconv.ParseBool(enableSwagger); err == nil {
			cfg.Server.EnableSwagger = b
		}
	}

	// Data config
	if dataDir := os.Getenv("LLAMADECK_DATA_DIRECTORY"); dataDir != "" {
		cfg.DataDir = dataDir
	}
	if logsDir := os.Getenv("LLAMADECK_LOGS_DIR"); logsDir != "" {
		cfg.Instances.LogsDir = logsDir
	}
	if instancesDir := os.Getenv("LLAMADECK_INSTANCES_DIR"); instancesDir != "" {
		cfg.Instances.InstancesDir = instancesDir
	}
	if autoCreate := os.Getenv("LLAMADECK_AUTO_CREATE_DATA_DIR"); autoCreate != "" {
		if b, err := strconv.ParseBool(autoCreate); err == nil {
			cfg.Instances.AutoCreateDirs = b
		}
	}

	// Instance config
	if portRange := os.Getenv("LLAMADECK_INSTANCE_PORT_RANGE"); portRange != "" {
		if ports := ParsePortRange(portRange); ports != [2]int{0, 0} {
			cfg.Instances.PortRange = ports
		}
	}
	if maxInstances := os.Getenv("LLAMADECK_MAX_INSTANCES"); maxInstances != "" {
		if m, err := strconv.Atoi(maxInstances); err == nil {
			cfg.Instances.MaxInstances = m
		}
	}
	if maxRunning := os.Getenv("LLAMADECK_MAX_RUNNING_INSTANCES"); maxRunning != "" {
		if m, err := strconv.Atoi(maxRunning); err == nil {
			cfg.Instances.MaxRunningInstances = m
		}
	}
	if enableLRUEviction := os.Getenv("LLAMADECK_ENABLE_LRU_EVICTION"); enableLRUEviction != "" {
		if b, err := strconv.ParseBool(enableLRUEviction); err == nil {
			cfg.Instances.EnableLRUEviction = b
		}
	}
	// Backend config
	// LlamaCpp backend
	if llamaCmd := os.Getenv("LLAMADECK_LLAMACPP_COMMAND"); llamaCmd != "" {
		cfg.Backends.LlamaCpp.Command = llamaCmd
	}
	if llamaArgs := os.Getenv("LLAMADECK_LLAMACPP_ARGS"); llamaArgs != "" {
		cfg.Backends.LlamaCpp.Args = strings.Split(llamaArgs, " ")
	}
	if llamaEnv := os.Getenv("LLAMADECK_LLAMACPP_ENV"); llamaEnv != "" {
		if cfg.Backends.LlamaCpp.Environment == nil {
			cfg.Backends.LlamaCpp.Environment = make(map[string]string)
		}
		parseEnvVars(llamaEnv, cfg.Backends.LlamaCpp.Environment)
	}
	if llamaDockerEnabled := os.Getenv("LLAMADECK_LLAMACPP_DOCKER_ENABLED"); llamaDockerEnabled != "" {
		if b, err := strconv.ParseBool(llamaDockerEnabled); err == nil {
			if cfg.Backends.LlamaCpp.Docker == nil {
				cfg.Backends.LlamaCpp.Docker = &DockerSettings{}
			}
			cfg.Backends.LlamaCpp.Docker.Enabled = b
		}
	}
	if llamaDockerImage := os.Getenv("LLAMADECK_LLAMACPP_DOCKER_IMAGE"); llamaDockerImage != "" {
		if cfg.Backends.LlamaCpp.Docker == nil {
			cfg.Backends.LlamaCpp.Docker = &DockerSettings{}
		}
		cfg.Backends.LlamaCpp.Docker.Image = llamaDockerImage
	}
	if llamaDockerArgs := os.Getenv("LLAMADECK_LLAMACPP_DOCKER_ARGS"); llamaDockerArgs != "" {
		if cfg.Backends.LlamaCpp.Docker == nil {
			cfg.Backends.LlamaCpp.Docker = &DockerSettings{}
		}
		cfg.Backends.LlamaCpp.Docker.Args = strings.Split(llamaDockerArgs, " ")
	}
	if llamaDockerEnv := os.Getenv("LLAMADECK_LLAMACPP_DOCKER_ENV"); llamaDockerEnv != "" {
		if cfg.Backends.LlamaCpp.Docker == nil {
			cfg.Backends.LlamaCpp.Docker = &DockerSettings{}
		}
		if cfg.Backends.LlamaCpp.Docker.Environment == nil {
			cfg.Backends.LlamaCpp.Docker.Environment = make(map[string]string)
		}
		parseEnvVars(llamaDockerEnv, cfg.Backends.LlamaCpp.Docker.Environment)
	}
	if llamaEnv := os.Getenv("LLAMADECK_LLAMACPP_RESPONSE_HEADERS"); llamaEnv != "" {
		if cfg.Backends.LlamaCpp.ResponseHeaders == nil {
			cfg.Backends.LlamaCpp.ResponseHeaders = make(map[string]string)
		}
		parseHeaders(llamaEnv, cfg.Backends.LlamaCpp.ResponseHeaders)
	}
	if llamaCacheDir := os.Getenv("LLAMADECK_LLAMACPP_CACHE_DIR"); llamaCacheDir != "" {
		cfg.Backends.LlamaCpp.CacheDir = llamaCacheDir
	}
	// Handle default llama.cpp env var
	if llamaCacheDir := os.Getenv("LLAMA_CACHE"); llamaCacheDir != "" {
		cfg.Backends.LlamaCpp.CacheDir = llamaCacheDir
	}
	if llamaDownloadTimeout := os.Getenv("LLAMADECK_LLAMACPP_DOWNLOAD_TIMEOUT"); llamaDownloadTimeout != "" {
		if t, err := strconv.Atoi(llamaDownloadTimeout); err == nil {
			cfg.Backends.LlamaCpp.DownloadTimeout = time.Duration(t) * time.Second
		}
	}

