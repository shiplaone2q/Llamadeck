package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type CachedModel struct {
	Node      string `json:"node"`
	Repo      string `json:"repo"`
	Tag       string `json:"tag,omitempty"`
	Files     []File `json:"files"`
	SizeBytes int64  `json:"size_bytes"`
}

type File struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	SizeBytes int64  `json:"size_bytes"`
	Type      string `json:"type"`
}

var (
	splitFilePattern = regexp.MustCompile(`-\d{5}-of-(\d{5})\.gguf$`)
	shaPattern       = regexp.MustCompile(`^[0-9a-f]{7,40}$`)
)

func ScanCache(cacheDir, nodeName string) ([]CachedModel, error) {
	models := []CachedModel{}

	fileManager := NewFileManager(cacheDir)
	hubRoot := fileManager.HFHubRoot()

	// Scan new HF snapshot layout
	if _, err := os.Stat(hubRoot); err == nil {
		entries, err := os.ReadDir(hubRoot)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "models--") {
				continue
			}
			dirName := entry.Name()
			repo := parseHFRepoDirName(dirName)
			if repo == "" {
				continue
			}

			// Iterate over all refs
			refsDir := filepath.Join(hubRoot, dirName, "refs")
			if refEntries, err := os.ReadDir(refsDir); err == nil {
				for _, refEntry := range refEntries {
					if refEntry.IsDir() {
						continue
					}
					commit, err := readRefFile(filepath.Join(refsDir, refEntry.Name()))
					if err != nil {
						continue
					}
					cachedModel := scanHFRepo(hubRoot, dirName, repo, commit, nodeName)
					if len(cachedModel.Files) > 0 {
						cachedModel.Tag = refEntry.Name()
						models = append(models, cachedModel)
					}
				}
			} else {
				// Fallback: if no refs, scan snapshots directory
				snapshotRoot := filepath.Join(hubRoot, dirName, "snapshots")
				if snapEntries, err := os.ReadDir(snapshotRoot); err == nil {
					for _, snapEntry := range snapEntries {
						if !snapEntry.IsDir() {
							continue
						}
						commit := snapEntry.Name()
						cachedModel := scanHFRepo(hubRoot, dirName, repo, commit, nodeName)
						if len(cachedModel.Files) > 0 {
							models = append(models, cachedModel)
