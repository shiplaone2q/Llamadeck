package server_test

import (
	"llamadeck/pkg/config"
	"llamadeck/pkg/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInferenceAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		managementKeys []string
		requestKey     string
		method         string
		expectedStatus int
	}{
		{
			name:           "valid management key for inference",
			managementKeys: []string{"sk-management-admin123"},
			requestKey:     "sk-management-admin123",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid key",
			managementKeys: []string{"sk-management-valid123"},
			requestKey:     "sk-management-invalid",
			method:         "GET",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing key",
			managementKeys: []string{"sk-management-valid123"},
			requestKey:     "",
			method:         "GET",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "OPTIONS request bypasses auth",
			managementKeys: []string{"sk-management-valid123"},
			requestKey:     "",
			method:         "OPTIONS",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "management key works for inference endpoint",
			managementKeys: []string{"sk-management-admin"},
			requestKey:     "sk-management-admin",
			method:         "POST",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.AuthConfig{
				RequireInferenceAuth: true,
				ManagementKeys:       tt.managementKeys,
			}
			middleware := server.NewAPIAuthMiddleware(cfg, nil)

			// Create test request
			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.requestKey != "" {
				req.Header.Set("Authorization", "Bearer "+tt.requestKey)
			}

			// Create test handler
