package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration with the following precedence:
// 1. Hardcoded defaults
// 2. .env file
// 3. Config file (with env var expansion)
// 4. LLAMADECK_* environment variable overrides
func LoadConfig(configPath string) (AppConfig, error) {
	// 1. Start with defaults
	defaultDataDir := getDefaultDataDir()
	cfg := getDefaultConfig(defaultDataDir)

	// 2. Load .env file
	loadDotEnv(configPath)

	// 3. Load from config file with env var expansion
	data, err := readConfigFile(configPath)
	if err != nil {
		return cfg, err
	}
	if data != nil {
		data = expandEnvVars(data)
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return cfg, err
		}
	}

	// If local node is not defined in nodes, add it with default config
	if _, ok := cfg.Nodes[cfg.LocalNode]; !ok {
		cfg.Nodes[cfg.LocalNode] = NodeConfig{}
	}

	// 4. Override with environment variables
	loadEnvVars(&cfg)

	if cfg.Auth.OIDC.Enabled() && !cfg.Auth.RequireManagementAuth {
		log.Printf("Warning: auth.oidc is enabled but require_management_auth is false; " +
