package backends_test

import (
	"llamadeck/pkg/backends"
	"llamadeck/pkg/testutil"
	"testing"
)

func TestParseVllmCommand(t *testing.T) {
	tests := []struct {
		name      string
		command   string
		expectErr bool
		validate  func(*testing.T, *backends.VllmServerOptions)
	}{
		{
			name:      "basic vllm serve command",
			command:   "vllm serve microsoft/DialoGPT-medium",
			expectErr: false,
			validate: func(t *testing.T, opts *backends.VllmServerOptions) {
				if opts.Model != "microsoft/DialoGPT-medium" {
					t.Errorf("expected model 'microsoft/DialoGPT-medium', got '%s'", opts.Model)
				}
			},
		},
		{
			name:      "serve only command",
			command:   "serve microsoft/DialoGPT-medium",
			expectErr: false,
			validate: func(t *testing.T, opts *backends.VllmServerOptions) {
				if opts.Model != "microsoft/DialoGPT-medium" {
					t.Errorf("expected model 'microsoft/DialoGPT-medium', got '%s'", opts.Model)
				}
			},
		},
		{
			name:      "positional model with flags",
			command:   "vllm serve microsoft/DialoGPT-medium --tensor-parallel-size 2",
			expectErr: false,
			validate: func(t *testing.T, opts *backends.VllmServerOptions) {
				if opts.Model != "microsoft/DialoGPT-medium" {
					t.Errorf("expected model 'microsoft/DialoGPT-medium', got '%s'", opts.Model)
				}
				if opts.TensorParallelSize != 2 {
					t.Errorf("expected tensor_parallel_size 2, got %d", opts.TensorParallelSize)
				}
			},
		},
		{
			name:      "model with path",
			command:   "vllm serve /path/to/model --gpu-memory-utilization 0.8",
			expectErr: false,
			validate: func(t *testing.T, opts *backends.VllmServerOptions) {
				if opts.Model != "/path/to/model" {
					t.Errorf("expected model '/path/to/model', got '%s'", opts.Model)
				}
				if opts.GPUMemoryUtilization != 0.8 {
					t.Errorf("expected gpu_memory_utilization 0.8, got %f", opts.GPUMemoryUtilization)
				}
			},
		},
		{
			name:      "multiple value types",
			command:   "vllm serve test-model --tensor-parallel-size 4 --gpu-memory-utilization 0.8 --enable-log-outputs",
			expectErr: false,
			validate: func(t *testing.T, opts *backends.VllmServerOptions) {
				if opts.Model != "test-model" {
					t.Errorf("expected model 'test-model', got '%s'", opts.Model)
				}
				if opts.TensorParallelSize != 4 {
					t.Errorf("expected tensor_parallel_size 4, got %d", opts.TensorParallelSize)
				}
				if opts.GPUMemoryUtilization != 0.8 {
					t.Errorf("expected gpu_memory_utilization 0.8, got %f", opts.GPUMemoryUtilization)
				}
				if !opts.EnableLogOutputs {
					t.Errorf("expected enable_log_outputs true, got %v", opts.EnableLogOutputs)
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
			command:   `vllm serve "unterminated`,
			expectErr: true,
		},
	}
