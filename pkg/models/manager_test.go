package models

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const testCommit = "abcdef1234567890abcdef1234567890abcdef12"

func sha256hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

type hfFixture struct {
	branches []HFBranch
	tree     []HFTreeEntry
	files    map[string][]byte
	etags    map[string]string
}

// newHFTestServer fakes the refs/tree/resolve endpoints and points HF_ENDPOINT at it.
func newHFTestServer(t *testing.T, fx hfFixture) *httptest.Server {
	t.Helper()

	if fx.branches == nil {
		fx.branches = []HFBranch{{Name: "main", TargetCommit: testCommit}}
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/api/models/") && strings.HasSuffix(r.URL.Path, "/refs"):
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(HFRefsResponse{Branches: fx.branches})
		case strings.HasPrefix(r.URL.Path, "/api/models/") && strings.Contains(r.URL.Path, "/tree/"):
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(fx.tree)
		case strings.Contains(r.URL.Path, "/resolve/"):
			// {repo}/resolve/{commit}/{path...}
			trimmed := strings.TrimPrefix(r.URL.Path, "/")
			idx := strings.Index(trimmed, "/resolve/")
			filePath := trimmed[idx+len("/resolve/"):]
			if slash := strings.Index(filePath, "/"); slash != -1 {
				filePath = filePath[slash+1:]
			}
			content, ok := fx.files[filePath]
			if !ok {
				http.NotFound(w, r)
				return
			}
			if etag := fx.etags[filePath]; etag != "" {
				w.Header().Set("ETag", etag)
			}
			http.ServeContent(w, r, filepath.Base(filePath), time.Time{}, bytes.NewReader(content))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	t.Setenv("HF_ENDPOINT", srv.URL)
	return srv
}

func newTestManager(t *testing.T) *Manager {
	t.Helper()
	cache := t.TempDir()
	t.Setenv("LLAMA_CACHE", cache)
	m := NewManager(cache, time.Minute, "test")
