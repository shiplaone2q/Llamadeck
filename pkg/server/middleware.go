package server

import (
	"context"
	"crypto/subtle"
	"fmt"
	"llamadeck/pkg/auth"
	"llamadeck/pkg/config"
	"llamadeck/pkg/database"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	apiKeyContextKey contextKey = "apiKey"
	userContextKey   contextKey = "user"
)

type APIAuthMiddleware struct {
	authStore             database.AuthStore
	requireInferenceAuth  bool
	requireManagementAuth bool
	managementKeys        map[string]bool // Config-based management keys
	sessions              *auth.SessionStore
}

// NewAPIAuthMiddleware creates a new APIAuthMiddleware with the given configuration
func NewAPIAuthMiddleware(authCfg config.AuthConfig, authStore database.AuthStore) *APIAuthMiddleware {
	// Load management keys from config into managementKeys map
	managementKeys := make(map[string]bool)
	for _, key := range authCfg.ManagementKeys {
		managementKeys[key] = true
	}

	// Handle legacy auto-generation for management keys if none provided and auth is required
	var generated bool = false
	const banner = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

	if authCfg.RequireManagementAuth && len(authCfg.ManagementKeys) == 0 {
		key, err := auth.GenerateKey("llamadeck-mgmt")
		if err != nil {
			log.Printf("Warning: Failed to generate management key: %v", err)
			// Fallback to PID-based key for safety
			key = fmt.Sprintf("sk-management-fallback-%d", os.Getpid())
		}
		managementKeys[key] = true
		generated = true
		fmt.Printf("%s\n⚠️  MANAGEMENT AUTHENTICATION REQUIRED\n%s\n", banner, banner)
		fmt.Printf("🔑  Generated Management API Key:\n\n    %s\n\n", key)
	}

	if generated {
		fmt.Printf("%s\n⚠️  IMPORTANT\n%s\n", banner, banner)
		fmt.Println("• This key is auto-generated and will change on restart")
		fmt.Println("• For production, add explicit keys to your configuration")
		fmt.Println("• Copy this key before it disappears from the terminal")
		fmt.Println(banner)
	}

	return &APIAuthMiddleware{
		authStore:             authStore,
		requireInferenceAuth:  authCfg.RequireInferenceAuth,
		requireManagementAuth: authCfg.RequireManagementAuth,
		managementKeys:        managementKeys,
