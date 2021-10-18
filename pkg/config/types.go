package config

import "time"

// BackendSettings contains structured backend configuration
type BackendSettings struct {
	Command         string            `yaml:"command" json:"command"`
	Args            []string          `yaml:"args" json:"args"`
	Environment     map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	Docker          *DockerSettings   `yaml:"docker,omitempty" json:"docker,omitempty"`
	ResponseHeaders map[string]string `yaml:"response_headers,omitempty" json:"response_headers,omitempty"`
	HealthPath      string            `yaml:"health_path,omitempty" json:"health_path,omitempty"`
	CacheDir        string            `yaml:"cache_dir,omitempty" json:"cache_dir,omitempty"`
	DownloadTimeout time.Duration     `yaml:"download_timeout,omitempty" json:"download_timeout,omitempty" swaggertype:"string" example:"3600s"`
}

// DockerSettings contains Docker-specific configuration
type DockerSettings struct {
	Enabled     bool              `yaml:"enabled" json:"enabled"`
	Image       string            `yaml:"image" json:"image"`
	Args        []string          `yaml:"args" json:"args"`
	Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
}

// BackendConfig contains backend executable configurations.
// Custom holds named user-defined backends (backends.custom.<name>).
type BackendConfig struct {
	LlamaCpp BackendSettings            `yaml:"llama-cpp" json:"llama-cpp"`
	VLLM     BackendSettings            `yaml:"vllm" json:"vllm"`
	MLX      BackendSettings            `yaml:"mlx" json:"mlx"`
	Custom   map[string]BackendSettings `yaml:"custom,omitempty" json:"custom,omitempty"`
}

// AppConfig represents the configuration for llamadeck
type AppConfig struct {
	Server    ServerConfig          `yaml:"server" json:"server"`
	Backends  BackendConfig         `yaml:"backends" json:"backends"`
	Instances InstancesConfig       `yaml:"instances" json:"instances"`
	Database  DatabaseConfig        `yaml:"database" json:"database"`
	Auth      AuthConfig            `yaml:"auth" json:"auth"`
	LocalNode string                `yaml:"local_node,omitempty" json:"local_node,omitempty"`
	Nodes     map[string]NodeConfig `yaml:"nodes,omitempty" json:"nodes,omitempty"`

	// Directory where all llamadeck data will be stored (database, instances, logs, etc.)
	DataDir string `yaml:"data_dir" json:"data_dir"`

	Version    string `yaml:"-" json:"version"`
	CommitHash string `yaml:"-" json:"commit_hash"`
	BuildTime  string `yaml:"-" json:"build_time"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	// Server host to bind to
