package server_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"llamadeck/pkg/config"
	"llamadeck/pkg/database"
	"llamadeck/pkg/manager"
	"llamadeck/pkg/server"
)

// newInstancesTestRouter wires the instance endpoints to a real manager over
// an in-memory DB, bypassing auth middleware.
func newInstancesTestRouter(t *testing.T) http.Handler {
	t.Helper()

	cfg := &config.AppConfig{
		Backends: config.BackendConfig{
			LlamaCpp: config.BackendSettings{
				Command: "sh",
				Args:    []string{"-c", "sleep 999999"},
			},
			MLX: config.BackendSettings{
				Command: "sh",
				Args:    []string{"-c", "sleep 999999"},
			},
		},
		Instances: config.InstancesConfig{
			PortRange:            [2]int{8000, 9000},
			MaxInstances:         10,
			MaxRunningInstances:  10,
			DefaultAutoRestart:   true,
			DefaultMaxRestarts:   3,
			LogsDir:              t.TempDir(),
			DefaultRestartDelay:  5,
			DefaultIdleTimeout:   30,
			TimeoutCheckInterval: 5,
		},
		Database: config.DatabaseConfig{
			Path:               ":memory:",
			MaxOpenConnections: 25,
			MaxIdleConnections: 5,
			ConnMaxLifetime:    5 * time.Minute,
		},
		LocalNode: "main",
		Nodes:     map[string]config.NodeConfig{},
	}

	db, err := database.Open(&database.Config{
		Path:               cfg.Database.Path,
		MaxOpenConnections: cfg.Database.MaxOpenConnections,
		MaxIdleConnections: cfg.Database.MaxIdleConnections,
