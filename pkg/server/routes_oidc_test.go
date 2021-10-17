package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"llamadeck/pkg/auth"
	"llamadeck/pkg/config"
)

// TestRoutesOIDCBypassManagementAuth verifies the router (not just the
// handlers) serves whoami and the OIDC endpoints without a management key.
func TestRoutesOIDCBypassManagementAuth(t *testing.T) {
	cfg := config.AppConfig{}
	cfg.Auth.RequireManagementAuth = true
	cfg.Auth.ManagementKeys = []string{"sk-test"}

	h := &Handler{cfg: cfg}
	h.authMiddleware = NewAPIAuthMiddleware(cfg.Auth, nil)
	h.oidc = &OIDCService{Sessions: auth.NewSessionStore(3600e9)}
	h.authMiddleware.sessions = h.oidc.Sessions

	r := SetupRouter(h)

