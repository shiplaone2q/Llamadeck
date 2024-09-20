package backends_test

import (
	"encoding/json"
	"llamadeck/pkg/backends"
	"llamadeck/pkg/config"
	"reflect"
	"testing"
)

func TestCustomBuildCommandArgs(t *testing.T) {
	tests := []struct {
		name     string
		options  backends.CustomServerOptions
		expected []string
	}{
		{
			name:     "port substitution",
			options:  backends.CustomServerOptions{Port: 8123, Args: []string{"serve", "--port", "{port}"}},
			expected: []string{"serve", "--port", "8123"},
		},
		{
			name:     "both placeholders in one arg",
			options:  backends.CustomServerOptions{Port: 9000, Model: "m1", Args: []string{"--run={model}-{port}"}},
			expected: []string{"--run=m1-9000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.options.BuildCommandArgs()
			if !reflect.DeepEqual(args, tt.expected) {
				t.Errorf("BuildCommandArgs() = %v, want %v", args, tt.expected)
			}
		})
	}
}

// Args are persisted; mutation would corrupt subsequent launches.
func TestCustomBuildCommandArgs_DoesNotMutateArgs(t *testing.T) {
	options := backends.CustomServerOptions{
		Port: 8123,
		Args: []string{"--port", "{port}"},
	}

	options.BuildCommandArgs()

	if options.Args[1] != "{port}" {
		t.Errorf("original args must not be mutated, got %v", options.Args)
	}
}

func TestOptionsGetHealthPath(t *testing.T) {
	backendConfig := &config.BackendConfig{
		Custom: map[string]config.BackendSettings{
			"x":    {HealthPath: "/ready"},
			"bare": {HealthPath: "ready"},
		},
		LlamaCpp: config.BackendSettings{HealthPath: "/llm-ready"},
	}

	tests := []struct {
		name     string
		options  backends.Options
		expected string
	}{
		{
			name:     "custom configured path",
			options:  backends.Options{BackendType: backends.BackendTypeCustom, CustomServerOptions: &backends.CustomServerOptions{Name: "x"}},
			expected: "/ready",
		},
		{
			name:     "missing leading slash is prefixed",
			options:  backends.Options{BackendType: backends.BackendTypeCustom, CustomServerOptions: &backends.CustomServerOptions{Name: "bare"}},
			expected: "/ready",
		},
		{
			name:     "entry without path falls back to default",
			options:  backends.Options{BackendType: backends.BackendTypeCustom, CustomServerOptions: &backends.CustomServerOptions{Name: "other"}},
			expected: "/health",
		},
		{
			name:     "built-in backends read the same field",
			options:  backends.Options{BackendType: backends.BackendTypeLlamaCpp, LlamaServerOptions: &backends.LlamaServerOptions{}},
			expected: "/llm-ready",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.options.GetHealthPath(backendConfig); got != tt.expected {
				t.Errorf("GetHealthPath() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestCustomValidate(t *testing.T) {
	tests := []struct {
