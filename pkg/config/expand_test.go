package config

import (
	"os"
	"testing"
)

func TestExpandEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "simple variable",
			input:    "host: ${MY_HOST}",
			envVars:  map[string]string{"MY_HOST": "localhost"},
			expected: "host: localhost",
		},
		{
			name:     "unset variable left as-is",
			input:    "host: ${UNSET_VAR}",
			envVars:  nil,
			expected: "host: ${UNSET_VAR}",
		},
		{
			name:     "variable with default when unset",
			input:    "port: ${MY_PORT:-8080}",
			envVars:  nil,
			expected: "port: 8080",
		},
		{
			name:     "variable with default when set",
			input:    "port: ${MY_PORT:-8080}",
			envVars:  map[string]string{"MY_PORT": "9090"},
			expected: "port: 9090",
		},
		{
			name:     "variable with empty default",
			input:    "val: ${MY_VAR:-}",
			envVars:  nil,
			expected: "val: ",
		},
		{
			name:     "variable with default when empty",
			input:    "val: ${MY_VAR:-fallback}",
			envVars:  map[string]string{"MY_VAR": ""},
			expected: "val: fallback",
		},
		{
			name:     "no braces not expanded",
			input:    "host: $MY_HOST",
			envVars:  map[string]string{"MY_HOST": "localhost"},
			expected: "host: $MY_HOST",
		},
		{
			name:     "multiple substitutions",
			input:    "${HOST}:${PORT}",
			envVars:  map[string]string{"HOST": "localhost", "PORT": "8080"},
			expected: "localhost:8080",
		},
		{
			name:     "inline in longer string",
