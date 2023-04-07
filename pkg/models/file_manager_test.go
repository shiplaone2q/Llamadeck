package models

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHFRepoDirName(t *testing.T) {
	fm := NewFileManager("/cache")

	tests := []struct {
		name     string
		repo     string
		expected string
	}{
		{
			name:     "simple repo",
			repo:     "bartowski/Llama-3.2-3B-GGUF",
			expected: "models--bartowski--Llama-3.2-3B-GGUF",
		},
		{
			name:     "nested repo",
			repo:     "bartowski/Qwen/Qwen3-30B-GGUF",
			expected: "models--bartowski--Qwen--Qwen3-30B-GGUF",
		},
		{
			name:     "org model",
			repo:     "org/model",
			expected: "models--org--model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fm.HFRepoDirName(tt.repo)
			if result != tt.expected {
				t.Errorf("HFRepoDirName(%q) = %q, want %q", tt.repo, result, tt.expected)
			}
		})
	}
}

func TestHFValidateBlobPath(t *testing.T) {
	cache := t.TempDir()
	t.Setenv("LLAMA_CACHE", cache)
	fm := NewFileManager(cache)

	blobDir := fm.HFRepoDir("org/model") + "/blobs"

	tests := []struct {
		name      string
		path      string
		wantValid bool
	}{
		{"hex blob under repo", blobDir + "/" + strings.Repeat("ab", 32), true},
