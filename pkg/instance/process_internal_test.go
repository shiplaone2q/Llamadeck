package instance

import (
	"llamadeck/pkg/backends"
	"llamadeck/pkg/config"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// Server returns 200 only on the custom path, 404 elsewhere.
func TestWaitForHealthy_CustomHealthPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ready":
			w.WriteHeader(http.StatusOK)
		case "/nocontent":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	host, portStr, err := net.SplitHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Fatalf("failed to parse test server address: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("failed to parse test server port: %v", err)
	}

	// The health path comes from backends.custom.<name>.health_path in
	// config, not from instance options.
	newTestInstance := func(healthPath string) *Instance {
		inst := &Instance{Name: "custom-health"}
		inst.status = newStatus(Running)
		inst.globalBackendSettings = &config.BackendConfig{
			Custom: map[string]config.BackendSettings{
				"test-backend": {HealthPath: healthPath},
			},
		}
		inst.options = newOptions(&Options{
			BackendOptions: backends.Options{
				BackendType: backends.BackendTypeCustom,
				CustomServerOptions: &backends.CustomServerOptions{
					Name: "test-backend",
					Host: host,
					Port: port,
				},
			},
		})
		inst.process = newProcess(inst)
		return inst
	}

	t.Run("configured path returns healthy", func(t *testing.T) {
		if err := newTestInstance("/ready").process.waitForHealthy(2); err != nil {
			t.Errorf("expected healthy on custom path, got: %v", err)
		}
	})

	t.Run("204 No Content counts as healthy", func(t *testing.T) {
