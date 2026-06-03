package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"llamadeck/pkg/backends"
	"llamadeck/pkg/instance"
	"llamadeck/pkg/validation"
	"log"
	"net/http"
	"strings"
)

// OpenAIListInstancesResponse represents the response structure for listing instances (models) in OpenAI-compatible format
type OpenAIListInstancesResponse struct {
	Object string           `json:"object"`
	Data   []OpenAIInstance `json:"data"`
}

// OpenAIInstance represents a single instance (model) in OpenAI-compatible format
type OpenAIInstance struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// LlamaCppModel represents a model available in a llama.cpp instance
type LlamaCppModel struct {
	ID      string              `json:"id"`
	Object  string              `json:"object"`
	OwnedBy string              `json:"owned_by"`
	Created int64               `json:"created"`
	InCache bool                `json:"in_cache"`
	Path    string              `json:"path"`
	Status  LlamaCppModelStatus `json:"status"`
}

// LlamaCppModelStatus represents the status of a model in a llama.cpp instance
type LlamaCppModelStatus struct {
	Value string   `json:"value"` // "loaded" | "loading" | "unloaded"
	Args  []string `json:"args"`
}

// fetchLlamaCppModels fetches models from a llama.cpp instance using the proxy
func fetchLlamaCppModels(inst *instance.Instance) ([]LlamaCppModel, error) {
	// Create a request to the instance's /models endpoint
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/models", inst.GetHost(), inst.GetPort()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Use a custom response writer to capture the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Data []LlamaCppModel `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data, nil
}

// OpenAIListInstances godoc
// @Summary List instances in OpenAI-compatible format
// @Description Returns a list of instances in a format compatible with OpenAI API
// @Tags OpenAI
// @Security ApiKeyAuth
// @Produces json
// @Success 200 {object} OpenAIListInstancesResponse "List of OpenAI-compatible instances"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/models [get]
func (h *Handler) OpenAIListInstances() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instances, err := h.InstanceManager.ListInstances()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "list_failed", "Failed to list instances: "+err.Error())
			return
		}

		var openaiInstances []OpenAIInstance

		// For each llama.cpp instance, try to fetch models and add them as separate entries
		for _, inst := range instances {

			if inst.GetBackendType() == backends.BackendTypeLlamaCpp && inst.IsRunning() {
