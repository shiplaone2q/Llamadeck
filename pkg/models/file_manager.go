package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	commitPattern = regexp.MustCompile(`^[0-9a-f]{40}$`)
)

type FileManager struct {
	cacheDir    string
	hubRootOnce sync.Once
	hubRootVal  string
}

func NewFileManager(cacheDir string) *FileManager {
	return &FileManager{cacheDir: cacheDir}
}

// sanitizePathComponent removes dangerous path elements to prevent path traversal.
func sanitizePathComponent(component string) string {
	component = strings.ReplaceAll(component, string(filepath.Separator), "_")
	component = strings.ReplaceAll(component, "/", "_")
	component = strings.ReplaceAll(component, "\\", "_")
	component = strings.ReplaceAll(component, "..", "_")
	component = strings.TrimPrefix(component, ".")
	component = strings.TrimSpace(component)
	if component == "" {
		component = "unknown"
	}
	return component
}

func (fm *FileManager) HFHubRoot() string {
	fm.hubRootOnce.Do(func() {
		for _, env := range []string{"LLAMA_CACHE", "HF_HUB_CACHE", "HUGGINGFACE_HUB_CACHE"} {
			if v := os.Getenv(env); v != "" {
				fm.hubRootVal = v
				return
			}
		}
		if v := os.Getenv("HF_HOME"); v != "" {
			fm.hubRootVal = filepath.Join(v, "hub")
			return
		}
		if v := os.Getenv("XDG_CACHE_HOME"); v != "" {
