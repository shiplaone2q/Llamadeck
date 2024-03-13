package backends_test

import (
	"llamadeck/pkg/backends"
	"llamadeck/pkg/config"
	"llamadeck/pkg/testutil"
	"testing"
)

func TestParseMlxCommand(t *testing.T) {
	tests := []struct {
		name      string
		command   string
		expectErr bool
		validate  func(*testing.T, *backends.MlxServerOptions)
	}{
		{
			name:      "basic command",
			command:   "mlx_lm.server --model /path/to/model --host 0.0.0.0",
			expectErr: false,
			validate: func(t *testing.T, opts *backends.MlxServerOptions) {
				if opts.Model != "/path/to/model" {
					t.Errorf("expected model '/path/to/model', got '%s'", opts.Model)
				}
				if opts.Host != "0.0.0.0" {
					t.Errorf("expected host '0.0.0.0', got '%s'", opts.Host)
				}
			},
		},
		{
			name:      "args only",
			command:   "--model /path/to/model --port 8080",
			expectErr: false,
			validate: func(t *testing.T, opts *backends.MlxServerOptions) {
				if opts.Model != "/path/to/model" {
					t.Errorf("expected model '/path/to/model', got '%s'", opts.Model)
				}
				if opts.Port != 8080 {
					t.Errorf("expected port 8080, got %d", opts.Port)
				}
			},
		},
		{
			name:      "mixed flag formats",
			command:   "mlx_lm.server --model=/path/model --temp=0.7 --trust-remote-code",
			expectErr: false,
			validate: func(t *testing.T, opts *backends.MlxServerOptions) {
				if opts.Model != "/path/model" {
					t.Errorf("expected model '/path/model', got '%s'", opts.Model)
				}
				if opts.Temp != 0.7 {
					t.Errorf("expected temp 0.7, got %f", opts.Temp)
				}
				if !opts.TrustRemoteCode {
					t.Errorf("expected trust_remote_code to be true")
				}
			},
		},
		{
			name:      "multiple value types",
			command:   "mlx_lm.server --model /test/model.mlx --port 8080 --temp 0.7 --trust-remote-code --log-level DEBUG",
			expectErr: false,
			validate: func(t *testing.T, opts *backends.MlxServerOptions) {
				if opts.Model != "/test/model.mlx" {
					t.Errorf("expected model '/test/model.mlx', got '%s'", opts.Model)
				}
				if opts.Port != 8080 {
					t.Errorf("expected port 8080, got %d", opts.Port)
				}
				if opts.Temp != 0.7 {
					t.Errorf("expected temp 0.7, got %f", opts.Temp)
				}
				if !opts.TrustRemoteCode {
					t.Errorf("expected trust_remote_code to be true")
				}
				if opts.LogLevel != "DEBUG" {
					t.Errorf("expected log_level 'DEBUG', got '%s'", opts.LogLevel)
				}
			},
		},
		{
			name:      "empty command",
			command:   "",
			expectErr: true,
		},
		{
			name:      "unterminated quote",
			command:   `mlx_lm.server --model test.mlx --chat-template "unterminated`,
			expectErr: true,
		},
