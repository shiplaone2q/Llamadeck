package config

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func loadDotEnv(explicitConfigPath string) string {
	var searchPaths []string

	if explicitConfigPath != "" {
		absPath, err := filepath.Abs(explicitConfigPath)
		if err == nil {
			searchPaths = append(searchPaths, filepath.Join(filepath.Dir(absPath), ".env"))
		}
	}

	cwd, err := os.Getwd()
	if err == nil {
		searchPaths = append(searchPaths, filepath.Join(cwd, ".env"))
	}

	for _, dir := range getPlatformConfigDirs() {
		searchPaths = append(searchPaths, filepath.Join(dir, ".env"))
	}

	for _, p := range searchPaths {
		data, err := os.ReadFile(p)
