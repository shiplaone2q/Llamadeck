package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"llamadeck/pkg/auth"
	"llamadeck/pkg/config"
)

// TestRoutesOIDCBypassManagementAuth verifies the router (not just the
