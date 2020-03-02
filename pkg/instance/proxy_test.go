package instance

import (
	"llamadeck/pkg/backends"
	"llamadeck/pkg/config"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestProxy_HealthRequestsExcluded(t *testing.T) {
	globalConfig := &config.AppConfig{
		Backends: config.BackendConfig{
			LlamaCpp: config.BackendSettings{
				Command:    "llama-server",
				HealthPath: "/health",
			},
		},
		Instances: config.InstancesConfig{
			LogsDir: t.TempDir(),
		},
		Nodes:     map[string]config.NodeConfig{},
		LocalNode: "main",
	}

	options := &Options{
		BackendOptions: backends.Options{
			BackendType: backends.BackendTypeLlamaCpp,
			LlamaServerOptions: &backends.LlamaServerOptions{
				Model: "/path/to/model.gguf",
				Port:  8080,
			},
		},
	}
