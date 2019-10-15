package server_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"llamadeck/pkg/config"
	"llamadeck/pkg/models"
	"llamadeck/pkg/server"
)

// newModelsTestHandler wires a real manager, temp cache, and dead HF endpoint.
func newModelsTestHandler(t *testing.T) (*server.Handler, *models.Manager, string) {
	t.Helper()
	cache := t.TempDir()
	t.Setenv("LLAMA_CACHE", cache)
	t.Setenv("HF_ENDPOINT", "http://127.0.0.1:1")
	mm := models.NewManager(cache, time.Second, "test")
	t.Cleanup(mm.Close)
