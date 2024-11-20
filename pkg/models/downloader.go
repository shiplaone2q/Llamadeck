package models

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func normalizeETag(etag string) string {
	etag = strings.TrimPrefix(etag, "W/")
	etag = strings.Trim(etag, "\"")
	return etag
}

const (
	progressChannelBuffer  = 100
	maxConcurrentDownloads = 5
	jobRetentionDuration   = 24 * time.Hour
	jobCleanupInterval     = 1 * time.Hour
	maxDownloadAttempts    = 3
	maxMetadataBytes       = 32 << 20
)

type ModelFormat string

const (
	FormatGGUF        ModelFormat = "gguf"
	FormatSafetensors ModelFormat = "safetensors"
)

func IsValidFormat(f ModelFormat) bool {
	return f == FormatGGUF || f == FormatSafetensors
}

const (
	hfAPIBase       = "https://huggingface.co"
	hfRefsURLFmt    = "%s/api/models/%s/refs"
	hfTreeURLFmt    = "%s/api/models/%s/tree/%s?recursive=true"
	hfResolveURLFmt = "%s/%s/resolve/%s/%s"
)

var lfsOIDPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

type HFBranch struct {
	Name         string `json:"name"`
	TargetCommit string `json:"targetCommit"`
}

type HFRefsResponse struct {
	Branches []HFBranch `json:"branches"`
	Tags     []HFBranch `json:"tags"`
}

type HFLFSInfo struct {
	OID  string `json:"oid"`
	Size int64  `json:"size"`
}

type HFTreeEntry struct {
	Path string     `json:"path"`
	Type string     `json:"type"`
	Size int64      `json:"size"`
	LFS  *HFLFSInfo `json:"lfs,omitempty"`
}

type HFRepoInfo struct {
	Repo   string        `json:"repo"`
	Commit string        `json:"commit"`
	Branch string        `json:"branch"`
	Files  []HFTreeEntry `json:"files"`
}

type HFDownloadTask struct {
	URL          string `json:"url"`
	BlobPath     string `json:"blob_path"`
	SnapshotPath string `json:"snapshot_path"`
	OID          string `json:"oid"`
	Filename     string `json:"filename"`
	Size         int64  `json:"size"`
	Type         string `json:"type"`
	LFS          bool   `json:"lfs"`
}

type HFDownloadPlan struct {
	Repo     HFRepoInfo       `json:"repo"`
	Tasks    []HFDownloadTask `json:"tasks"`
	MainGGUF *HFDownloadTask  `json:"main_gguf"`
	MMProj   *HFDownloadTask  `json:"mmproj"`
	Preset   *HFDownloadTask  `json:"preset"`
	Format   ModelFormat      `json:"format"`
}

type Downloader struct {
	httpClient      *http.Client
	cacheDir        string
	version         string
	fileManager     *FileManager
	progressTracker *ProgressTracker
}

func NewDownloader(cacheDir string, timeout time.Duration, version string, fileManager *FileManager, progressTracker *ProgressTracker) *Downloader {
	if timeout == 0 {
		timeout = 60 * time.Minute
	}
	if version == "" {
		version = "0.1.0"
	}
	return &Downloader{
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:          10,
				IdleConnTimeout:       30 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 30 * time.Second,
				DisableCompression:    true,
			},
		},
		cacheDir:        cacheDir,
		version:         version,
		fileManager:     fileManager,
		progressTracker: progressTracker,
	}
}

func (md *Downloader) FetchRefs(ctx context.Context, repo, branch string) (string, error) {
	if branch == "" {
		branch = "main"
	}

	baseURL := md.hfBaseURL()
	url := fmt.Sprintf(hfRefsURLFmt, baseURL, escapeHFPath(repo))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	md.setCommonHeaders(req)

	resp, err := md.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch refs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("model not found: %s", repo)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return "", fmt.Errorf("failed to fetch refs: HTTP %d", resp.StatusCode)
	}

	var refsResp HFRefsResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxMetadataBytes)).Decode(&refsResp); err != nil {
		return "", fmt.Errorf("failed to decode refs response: %w", err)
	}

	// Find the requested branch by name, fall back to first branch
	for _, b := range refsResp.Branches {
		if b.Name == branch {
			return validatedCommit(b.TargetCommit, b.Name)
		}
	}
	if len(refsResp.Branches) > 0 {
		first := refsResp.Branches[0]
		return validatedCommit(first.TargetCommit, first.Name)
	}

