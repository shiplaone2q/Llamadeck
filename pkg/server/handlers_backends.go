package server

import (
	"encoding/json"
	"fmt"
	"llamadeck/pkg/backends"
	"llamadeck/pkg/instance"
	"net/http"
	"os/exec"
	"strings"
)

// ParseCommandRequest represents the request body for backend command parsing
type ParseCommandRequest struct {
	Command string `json:"command"`
}

// validateLlamaCppInstance validates that the instance specified in the request is a llama.cpp instance
func (h *Handler) validateLlamaCppInstance(r *http.Request) (*instance.Instance, error) {
	inst, err := h.getInstance(r)
	if err != nil {
		return nil, fmt.Errorf("invalid instance: %w", err)
	}

	options := inst.GetOptions()
	if options == nil {
		return nil, fmt.Errorf("cannot obtain instance's options")
	}

	if options.BackendOptions.BackendType != backends.BackendTypeLlamaCpp {
		return nil, fmt.Errorf("instance is not a llama.cpp server")
	}

	return inst, nil
}

// stripLlamaCppPrefix removes the llama.cpp proxy prefix from the request URL path
func (h *Handler) stripLlamaCppPrefix(r *http.Request, instName string) {
	// Strip the "/llama-cpp/<name>" prefix from the request URL
	prefix := fmt.Sprintf("/llama-cpp/%s", instName)
	r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
}

// LlamaCppUIProxy godoc
// @Summary Proxy requests to llama.cpp UI for the instance
// @Description Proxies requests to the llama.cpp UI for the specified instance
// @Tags Llama.cpp
// @Security ApiKeyAuth
// @Produce html
// @Param name query string true "Instance Name"
// @Success 200 {string} string "Proxied HTML response"
// @Failure 400 {string} string "Invalid instance"
// @Failure 404 {string} string "Instance does not exist"
// @Failure 500 {string} string "Internal Server Error"
// @Router /llama-cpp/{name}/ [get]
func (h *Handler) LlamaCppUIProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		inst, err := h.validateLlamaCppInstance(r)
		if err != nil {
			writeInstanceLookupError(w, err)
			return
		}

		if !inst.IsRemote() && !inst.IsRunning() {
			writeError(w, http.StatusBadRequest, "instance is not running", "Instance is not running")
			return
		}

		if !inst.IsRemote() {
			h.stripLlamaCppPrefix(r, inst.Name)
		}

		// Use instance's ServeHTTP which tracks inflight requests and handles shutting down state
		err = inst.ServeHTTP(w, r)
		if err != nil {
			// Error is already handled in ServeHTTP (response written)
			return
		}
	}
}

// LlamaCppProxy godoc
// @Summary Proxy requests to llama.cpp server instance
// @Description Proxies requests to the specified llama.cpp server instance, starting it on-demand if configured
// @Tags Llama.cpp
// @Security ApiKeyAuth
// @Produce json
// @Param name path string true "Instance Name"
// @Success 200 {object} map[string]any "Proxied response"
// @Failure 400 {string} string "Invalid instance"
// @Failure 404 {string} string "Instance does not exist"
// @Failure 500 {string} string "Internal Server Error"
// @Router /llama-cpp/{name}/props [get]
// @Router /llama-cpp/{name}/slots [get]
// @Router /llama-cpp/{name}/apply-template [post]
// @Router /llama-cpp/{name}/completion [post]
// @Router /llama-cpp/{name}/detokenize [post]
// @Router /llama-cpp/{name}/embeddings [post]
// @Router /llama-cpp/{name}/infill [post]
// @Router /llama-cpp/{name}/metrics [post]
// @Router /llama-cpp/{name}/props [post]
// @Router /llama-cpp/{name}/reranking [post]
// @Router /llama-cpp/{name}/tokenize [post]
func (h *Handler) LlamaCppProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		inst, err := h.validateLlamaCppInstance(r)
		if err != nil {
			writeInstanceLookupError(w, err)
			return
		}

		// Check instance permissions and resolve lifecycle flags in one query
		perm, err := h.authMiddleware.ResolveInstancePermission(r.Context(), inst.ID)
		if err != nil {
			writeError(w, http.StatusForbidden, "permission_denied", err.Error())
			return
		}

		// Check if instance is shutting down before autostart logic
		if inst.GetStatus() == instance.ShuttingDown {
			writeError(w, http.StatusServiceUnavailable, "instance_shutting_down", "Instance is shutting down")
			return
		}

		if !inst.IsRemote() && !inst.IsRunning() {
			if err := h.ensureInstanceRunning(inst, perm.CanStart, perm.CanEvict); err != nil {
				writeStartError(w, err)
				return
			}
		}

		if !inst.IsRemote() {
			h.stripLlamaCppPrefix(r, inst.Name)
		}

		// Use instance's ServeHTTP which tracks inflight requests and handles shutting down state
		err = inst.ServeHTTP(w, r)
		if err != nil {
			// Error is already handled in ServeHTTP (response written)
			return
		}
	}
}

// parseHelper parses a backend command and returns the parsed options
func parseHelper(w http.ResponseWriter, r *http.Request, backend interface {
