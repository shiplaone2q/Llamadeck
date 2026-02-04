package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"

	"llamadeck/pkg/config"
)

// NodeResponse represents a node configuration in API responses
type NodeResponse struct {
	Address string `json:"address"`
}

// NodeHealthStatus represents the health status of a node
type NodeHealthStatus string

const (
	NodeStatusHealthy     NodeHealthStatus = "healthy"
	NodeStatusUnreachable NodeHealthStatus = "unreachable"
)

// NodeHealthResponse represents the health status of a node in API responses
type NodeHealthResponse struct {
	Status    NodeHealthStatus `json:"status"`
	LatencyMS int64            `json:"latency_ms"` // no omitempty: 0 is meaningful (local node, sub-ms ping)
	Error     string           `json:"error,omitempty"`
	CheckedAt int64            `json:"checked_at"`
}

// ListNodes godoc
// @Summary List all configured nodes
// @Description Returns a map of all nodes configured in the server (node name -> node config)
// @Tags Nodes
// @Security ApiKeyAuth
// @Produces json
// @Success 200 {object} map[string]NodeResponse "Map of nodes"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/nodes [get]
func (h *Handler) ListNodes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Convert to sanitized response format (map of name -> NodeResponse)
		nodeResponses := make(map[string]NodeResponse, len(h.cfg.Nodes))
		for name, node := range h.cfg.Nodes {
			nodeResponses[name] = NodeResponse{
				Address: node.Address,
			}
		}

		writeJSON(w, http.StatusOK, nodeResponses)
	}
}

// GetNode godoc
// @Summary Get details of a specific node
// @Description Returns the details of a specific node by name
// @Tags Nodes
// @Security ApiKeyAuth
// @Produces json
// @Param name path string true "Node Name"
// @Success 200 {object} NodeResponse "Node details"
