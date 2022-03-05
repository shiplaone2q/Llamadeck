package models

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var repoNamePattern = regexp.MustCompile(
	`^[A-Za-z0-9][A-Za-z0-9._-]*(/[A-Za-z0-9][A-Za-z0-9._-]*)+$`)

type Manager struct {
	cacheDir    string
	jobStore    *JobStore
	downloader  *Downloader
	fileManager *FileManager
}

func NewManager(cacheDir string, timeout time.Duration, version string) *Manager {
	jobStore := NewJobStore()
	fileManager := NewFileManager(cacheDir)
	progressTracker := NewProgressTracker(jobStore)
	downloader := NewDownloader(cacheDir, timeout, version, fileManager, progressTracker)

	m := &Manager{
		cacheDir:    cacheDir,
		jobStore:    jobStore,
		downloader:  downloader,
		fileManager: fileManager,
	}
	// no download can be live at construction, so leftover scratch files are stale
	m.sweepStaleTempFiles()
	return m
}

func (m *Manager) sweepStaleTempFiles() {
	_ = filepath.WalkDir(m.fileManager.HFHubRoot(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("stale temp sweep: %s: %v", path, err)
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		name := d.Name()
		// dot prefix = scratch marker; legit names can't start with "." (sanitizePathComponent strips it)
		if !strings.HasPrefix(name, ".") ||
			(!strings.Contains(name, ".tmp-") && !strings.Contains(name, ".lnk-")) {
			return nil
		}
		if err := os.Remove(path); err != nil {
			log.Printf("stale temp sweep: failed to remove %s: %v", path, err)
			return nil
		}
		log.Printf("stale temp sweep: removed %s", path)
		return nil
	})
}

func (m *Manager) StartDownload(repo, tag string, format ModelFormat) (string, error) {
	if format == "" {
