package backends_test

import (
	"encoding/json"
	"fmt"
	"llamadeck/pkg/backends"
	"llamadeck/pkg/config"
	"llamadeck/pkg/testutil"
	"reflect"
	"testing"
)

func TestLlamaCppBuildCommandArgs_BooleanFields(t *testing.T) {
	tests := []struct {
		name     string
		options  backends.LlamaServerOptions
		expected []string
		excluded []string
	}{
		{
			name: "verbose true",
			options: backends.LlamaServerOptions{
				Verbose: true,
			},
			expected: []string{"--verbose"},
		},
		{
			name: "verbose false",
			options: backends.LlamaServerOptions{
				Verbose: false,
			},
			excluded: []string{"--verbose"},
		},
		{
			name: "multiple booleans",
			options: backends.LlamaServerOptions{
				Verbose: true,
				Mlock:   false,
				NoMmap:  true,
			},
			expected: []string{"--verbose", "--no-mmap"},
			excluded: []string{"--mlock"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.options.BuildCommandArgs()

			for _, expectedArg := range tt.expected {
				if !testutil.Contains(args, expectedArg) {
					t.Errorf("Expected argument %q not found in %v", expectedArg, args)
				}
			}

			for _, excludedArg := range tt.excluded {
				if testutil.Contains(args, excludedArg) {
					t.Errorf("Excluded argument %q found in %v", excludedArg, args)
				}
			}
		})
	}
}

func TestLlamaCppBuildCommandArgs_ZeroValues(t *testing.T) {
	options := backends.LlamaServerOptions{
		Port:        0,     // Should be excluded
		Threads:     0,     // Should be excluded
		Temperature: 0,     // Should be excluded
		Model:       "",    // Should be excluded
		Verbose:     false, // Should be excluded
	}

	args := options.BuildCommandArgs()

	// Zero values should not appear in arguments
	excludedArgs := []string{
		"--port", "0",
		"--threads", "0",
		"--temperature", "0",
		"--model", "",
		"--verbose",
	}

	for _, excludedArg := range excludedArgs {
		if testutil.Contains(args, excludedArg) {
			t.Errorf("Zero value argument %q should not be present in %v", excludedArg, args)
		}
	}
}

func TestLlamaCppBuildCommandArgs_ArrayFields(t *testing.T) {
	options := backends.LlamaServerOptions{
		Lora:               []string{"adapter1.bin", "adapter2.bin"},
		OverrideTensor:     []string{"tensor1", "tensor2", "tensor3"},
		DrySequenceBreaker: []string{".", "!", "?"},
	}

	args := options.BuildCommandArgs()

	// Check that each array value appears with its flag
	expectedOccurrences := map[string][]string{
		"--lora":                 {"adapter1.bin", "adapter2.bin"},
		"--override-tensor":      {"tensor1", "tensor2", "tensor3"},
		"--dry-sequence-breaker": {".", "!", "?"},
	}

	for flag, values := range expectedOccurrences {
		for _, value := range values {
			if !testutil.ContainsFlagWithValue(args, flag, value) {
				t.Errorf("Expected %s %s, not found in %v", flag, value, args)
			}
		}
	}
}

func TestLlamaCppBuildCommandArgs_EmptyArrays(t *testing.T) {
	options := backends.LlamaServerOptions{
		Lora:           []string{}, // Empty array should not generate args
		OverrideTensor: []string{}, // Empty array should not generate args
	}

	args := options.BuildCommandArgs()

	excludedArgs := []string{"--lora", "--override-tensor"}
	for _, excludedArg := range excludedArgs {
		if testutil.Contains(args, excludedArg) {
			t.Errorf("Empty array should not generate argument %q in %v", excludedArg, args)
		}
	}
}

func TestLlamaCppUnmarshalJSON_StandardFields(t *testing.T) {
	jsonData := `{
		"model": "/path/to/model.gguf",
		"port": 8080,
		"host": "localhost", 
		"verbose": true,
		"ctx_size": 4096,
		"gpu_layers": 32,
		"temp": 0.7
	}`

	var options backends.LlamaServerOptions
	err := json.Unmarshal([]byte(jsonData), &options)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if options.Model != "/path/to/model.gguf" {
		t.Errorf("Expected model '/path/to/model.gguf', got %q", options.Model)
	}
	if options.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", options.Port)
	}
	if options.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got %q", options.Host)
	}
	if !options.Verbose {
		t.Error("Expected verbose to be true")
	}
	if options.CtxSize != 4096 {
		t.Errorf("Expected ctx_size 4096, got %d", options.CtxSize)
	}
	if options.GPULayers != 32 {
		t.Errorf("Expected gpu_layers 32, got %d", options.GPULayers)
	}
	if options.Temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", options.Temperature)
	}
}

func TestLlamaCppUnmarshalJSON_AlternativeFieldNames(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		checkFn  func(backends.LlamaServerOptions) error
	}{
		{
			name:     "threads alternatives",
			jsonData: `{"t": 4, "tb": 2}`,
			checkFn: func(opts backends.LlamaServerOptions) error {
				if opts.Threads != 4 {
					return fmt.Errorf("expected threads 4, got %d", opts.Threads)
				}
				if opts.ThreadsBatch != 2 {
					return fmt.Errorf("expected threads_batch 2, got %d", opts.ThreadsBatch)
				}
				return nil
			},
		},
		{
			name:     "context size alternatives",
			jsonData: `{"c": 2048}`,
			checkFn: func(opts backends.LlamaServerOptions) error {
				if opts.CtxSize != 2048 {
					return fmt.Errorf("expected ctx_size 4096, got %d", opts.CtxSize)
