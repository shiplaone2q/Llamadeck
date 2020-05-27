package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"llamadeck/pkg/config"
	"llamadeck/pkg/server"
)

func newNodesTestHandler(t *testing.T, nodes map[string]config.NodeConfig, localNode string) http.Handler {
	t.Helper()
	cfg := config.AppConfig{
		LocalNode: localNode,
		Nodes:     nodes,
	}
	h := server.NewHandler(nil, nil, cfg, nil, nil)
	return h.GetNodeHealth()
}

func TestGetNodeHealth(t *testing.T) {
	remote := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/version" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("Authorization") != "Bearer node-key" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer remote.Close()

	tests := []struct {
		name           string
		nodes          map[string]config.NodeConfig
		localNode      string
		requestName    string
		expectedStatus int
		expectedHealth server.NodeHealthStatus
	}{
		{
			name:           "local node is healthy without network call",
