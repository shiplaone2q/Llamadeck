package server

import (
	"encoding/json"
	"fmt"
	"llamadeck/pkg/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// CreateKeyRequest represents the request body for creating a new API key.
type CreateKeyRequest struct {
	Name           string                   `json:"name"`
	PermissionMode auth.PermissionMode      `json:"permission_mode"`
	ExpiresAt      *int64                   `json:"expires_at,omitempty"`
	Permissions    []InstancePermissionSpec `json:"permissions,omitempty"`
}

// InstancePermissionSpec specifies per-instance access for a new key. CanStart and
// CanEvict default to true when omitted or null.
type InstancePermissionSpec struct {
	InstanceID int   `json:"instance_id"`
	CanStart   *bool `json:"can_start,omitempty"`
	CanEvict   *bool `json:"can_evict,omitempty"`
}

// CreateKeyResponse represents the response returned when creating a new API key.
type CreateKeyResponse struct {
	ID             int                 `json:"id"`
	Name           string              `json:"name"`
	UserID         string              `json:"user_id"`
	PermissionMode auth.PermissionMode `json:"permission_mode"`
	ExpiresAt      *int64              `json:"expires_at"`
	CreatedAt      int64               `json:"created_at"`
	UpdatedAt      int64               `json:"updated_at"`
	LastUsedAt     *int64              `json:"last_used_at"`
	Key            string              `json:"key"`
}

// KeyResponse represents an API key in responses for list and get operations.
type KeyResponse struct {
	ID             int                 `json:"id"`
	Name           string              `json:"name"`
	UserID         string              `json:"user_id"`
	PermissionMode auth.PermissionMode `json:"permission_mode"`
	ExpiresAt      *int64              `json:"expires_at"`
	CreatedAt      int64               `json:"created_at"`
	UpdatedAt      int64               `json:"updated_at"`
	LastUsedAt     *int64              `json:"last_used_at"`
}

// KeyPermissionResponse represents the permissions for an API key on a specific instance.
type KeyPermissionResponse struct {
	InstanceID   int    `json:"instance_id"`
	InstanceName string `json:"instance_name"`
	CanStart     bool   `json:"can_start"`
	CanEvict     bool   `json:"can_evict"`
}

// CreateKey godoc
// @Summary Create a new API key
// @Description Creates a new API key with the specified permissions and returns the plain-text key (only shown once)
// @Tags Keys
// @Accept json
// @Produce json
// @Param key body CreateKeyRequest true "API key configuration"
// @Success 201 {object} CreateKeyResponse "Created API key with plain-text key"
// @Failure 400 {string} string "Invalid request body or validation error"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/auth/keys [post]
func (h *Handler) CreateKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateKeyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON in request body")
			return
		}

		// Validate request
		if req.Name == "" {
			writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
			return
		}
		if len(req.Name) > 100 {
			writeError(w, http.StatusBadRequest, "invalid_name", "Name must be 100 characters or less")
			return
		}
		if req.PermissionMode != auth.PermissionModeAllowAll && req.PermissionMode != auth.PermissionModePerInstance {
