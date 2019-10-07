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
