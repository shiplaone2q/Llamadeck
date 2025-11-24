package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"llamadeck/pkg/config"
	"llamadeck/pkg/database"
	"llamadeck/pkg/instance"
	"llamadeck/pkg/manager"
	"llamadeck/pkg/models"
	"llamadeck/pkg/validation"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

// Sentinel errors for instance start outcomes. Used by callers (via errors.Is)
// to classify the failure into the right HTTP response instead of string matching.
var (
	// ErrInstanceNotRunning is returned when the key is not permitted to start
	// the instance (can_start=false).
	ErrInstanceNotRunning = errors.New("instance is not running")
	// ErrMaxInstancesReached is returned when capacity is full and the key cannot
	// (or may not) evict to make room.
	ErrMaxInstancesReached = errors.New("cannot start instance, maximum number of instances reached")
)

// errorResponse represents an error response returned by the API
type errorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// writeError writes a JSON error response with the specified HTTP status code
func writeError(w http.ResponseWriter, status int, code, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(errorResponse{Error: code, Details: details}); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

// writeJSON writes a JSON response with the specified HTTP status code
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// writeText writes a plain text response with the specified HTTP status code
func writeText(w http.ResponseWriter, status int, data string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	if _, err := w.Write([]byte(data)); err != nil {
		log.Printf("Failed to write text response: %v", err)
	}
}

// Handler provides HTTP handlers for the llamadeck server API
type Handler struct {
	InstanceManager manager.InstanceManager
	modelManager    *models.Manager
	cfg             config.AppConfig
	httpClient      *http.Client
	authStore       database.AuthStore
	authMiddleware  *APIAuthMiddleware
