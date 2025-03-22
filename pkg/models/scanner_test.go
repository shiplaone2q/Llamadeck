package models

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseHFRepoDirName(t *testing.T) {
	tests := []struct {
		name     string
		dirName  string
		wantRepo string
	}{
		{
			name:     "simple repo",
			dirName:  "models--bartowski--Llama-3.2-3B-GGUF",
			wantRepo: "bartowski/Llama-3.2-3B-GGUF",
		},
		{
			name:     "nested repo",
			dirName:  "models--bartowski--Qwen--Qwen3-30B-GGUF",
			wantRepo: "bartowski/Qwen/Qwen3-30B-GGUF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := parseHFRepoDirName(tt.dirName)
			if repo != tt.wantRepo {
				t.Errorf("parseHFRepoDirName(%q) = %q, want %q", tt.dirName, repo, tt.wantRepo)
			}
		})
	}
}

func TestScanCache(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("LLAMA_CACHE", tmpDir)

	fm := NewFileManager(tmpDir)
	hubRoot := fm.HFHubRoot()

	repoDir := filepath.Join(hubRoot, "models--bartowski--Llama-3.2-3B-GGUF")
	blobsDir := filepath.Join(repoDir, "blobs")
	snapshotDir := filepath.Join(repoDir, "snapshots", "abc123def456789012345678901234567890abcd")
	refsDir := filepath.Join(repoDir, "refs")

	os.MkdirAll(blobsDir, 0755)
	os.MkdirAll(snapshotDir, 0755)
	os.MkdirAll(refsDir, 0755)

	os.WriteFile(filepath.Join(blobsDir, "blob1"), []byte("fake model data"), 0644)
	os.WriteFile(filepath.Join(blobsDir, "blob2"), []byte("fake mmproj data"), 0644)
	os.WriteFile(filepath.Join(refsDir, "main"), []byte("abc123def456789012345678901234567890abcd"), 0644)

	os.Symlink("../../blobs/blob1", filepath.Join(snapshotDir, "model-Q4_K_M.gguf"))
	os.Symlink("../../blobs/blob2", filepath.Join(snapshotDir, "mmproj-F16.gguf"))

	models, err := ScanCache(tmpDir, "test-node")
	if err != nil {
		t.Fatalf("ScanCache failed: %v", err)
	}

	if len(models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(models))
	}

	model := models[0]
	if model.Node != "test-node" {
		t.Errorf("node = %q, want %q", model.Node, "test-node")
	}

	if model.Repo != "bartowski/Llama-3.2-3B-GGUF" {
		t.Errorf("repo = %q, want %q", model.Repo, "bartowski/Llama-3.2-3B-GGUF")
	}

	if len(model.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(model.Files))
	}

	hasGGUF := false
	hasMMProj := false
	for _, file := range model.Files {
		if file.Type == "gguf" {
			hasGGUF = true
		}
		if file.Type == "mmproj" {
			hasMMProj = true
		}
	}

