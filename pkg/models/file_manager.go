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
