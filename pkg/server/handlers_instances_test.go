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
