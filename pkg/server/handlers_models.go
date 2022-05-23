package server

import (
	"encoding/json"
	"fmt"
	"llamadeck/pkg/config"
	"llamadeck/pkg/models"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
)

// DownloadRequest represents the request body for initiating a model download
type DownloadRequest struct {
	Repo   string             `json:"repo"`
	Format models.ModelFormat `json:"format"`
}

// DownloadResponse represents the response after initiating a model download
type DownloadResponse struct {
	JobID string `json:"job_id"`
	Repo  string `json:"repo"`
	Tag   string `json:"tag"`
}

// JobResponse represents the details of a download job for API responses
type JobResponse struct {
	ID          string          `json:"id"`
	Repo        string          `json:"repo"`
	Tag         string          `json:"tag"`
	Status      string          `json:"status"`
	Progress    models.Progress `json:"progress"`
	Error       string          `json:"error,omitempty"`
	CreatedAt   int64           `json:"created_at"`
	CompletedAt *int64          `json:"completed_at,omitempty"`
	ModelPath   string          `json:"model_path,omitempty"`
}

// ListJobsResponse represents the response for listing all download jobs
type ListJobsResponse struct {
	Jobs []JobResponse `json:"jobs"`
}

// shouldForwardToNode checks if the request should be forwarded to a remote node
func (h *Handler) shouldForwardToNode(nodeName string) bool {
	return nodeName != "" && nodeName != h.cfg.LocalNode
}

// forwardToNode forwards an HTTP request to a remote node
func (h *Handler) forwardToNode(nodeName string, w http.ResponseWriter, r *http.Request) bool {
	node, exists := h.cfg.Nodes[nodeName]
	if !exists {
		writeError(w, http.StatusNotFound, "node_not_found", "Node not found")
		return false
	}

	targetURL, err := url.Parse(node.Address)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "invalid_node_address", "Failed to parse node address")
		return false
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Transport = &authTransport{
		base:   http.DefaultTransport,
		apiKey: node.APIKey,
	}

	proxy.ServeHTTP(w, r)
	return true
}

type authTransport struct {
	base   http.RoundTripper
	apiKey string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+t.apiKey)
	}
	return t.base.RoundTrip(req)
}

// DownloadModel godoc
// @Summary Download a model from a repository
// @Description Initiates the download of a model from a specified repository and tag. Returns a job ID to track progress.
// @Tags Models
// @Security ApiKeyAuth
// @Accept json
// @Produce json
