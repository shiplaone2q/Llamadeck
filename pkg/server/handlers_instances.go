package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"llamadeck/pkg/instance"
	"llamadeck/pkg/manager"
	"llamadeck/pkg/validation"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// ListInstances godoc
// @Summary List all instances
// @Description Returns a list of all instances managed by the server
// @Tags Instances
// @Security ApiKeyAuth
// @Produces json
// @Success 200 {array} instance.Instance "List of instances"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/instances [get]
func (h *Handler) ListInstances() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instances, err := h.InstanceManager.ListInstances()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "list_failed", "Failed to list instances: "+err.Error())
			return
		}

		writeJSON(w, http.StatusOK, instances)
	}
}

// CreateInstance godoc
// @Summary Create and start a new instance
// @Description Creates a new instance with the provided configuration options
// @Tags Instances
// @Security ApiKeyAuth
// @Accept json
// @Produces json
// @Param name path string true "Instance Name"
// @Param options body instance.Options true "Instance configuration options"
// @Success 201 {object} instance.Instance "Created instance details"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/instances/{name} [post]
func (h *Handler) CreateInstance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		validatedName, err := validation.ValidateInstanceName(name)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_instance_name", err.Error())
			return
		}

		var options instance.Options
		if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
			return
		}

		inst, err := h.InstanceManager.CreateInstance(validatedName, &options)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "create_failed", "Failed to create instance: "+err.Error())
			return
		}

		writeJSON(w, http.StatusCreated, inst)
	}
}

// GetInstance godoc
// @Summary Get details of a specific instance
// @Description Returns the details of a specific instance by name. Returns
// @Description 404 when the instance does not exist locally or on its node.
// @Tags Instances
// @Security ApiKeyAuth
// @Produces json
// @Param name path string true "Instance Name"
// @Success 200 {object} instance.Instance "Instance details"
// @Failure 400 {string} string "Invalid name format"
// @Failure 404 {string} string "Instance does not exist"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/instances/{name} [get]
func (h *Handler) GetInstance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		validatedName, err := validation.ValidateInstanceName(name)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_instance_name", err.Error())
			return
		}

		inst, err := h.InstanceManager.GetInstance(validatedName)
		if err != nil {
			if errors.Is(err, manager.ErrInstanceNotFound) {
				writeError(w, http.StatusNotFound, "not_found", err.Error())
				return
			}
			writeError(w, http.StatusInternalServerError, "get_failed", "Failed to get instance: "+err.Error())
			return
		}

		writeJSON(w, http.StatusOK, inst)
	}
}

// UpdateInstance godoc
// @Summary Update an instance's configuration
// @Description Updates the configuration of a specific instance by name
// @Tags Instances
// @Security ApiKeyAuth
// @Accept json
// @Produces json
// @Param name path string true "Instance Name"
// @Param options body instance.Options true "Instance configuration options"
// @Success 200 {object} instance.Instance "Updated instance details"
// @Failure 400 {string} string "Invalid name format"
// @Failure 404 {string} string "Instance does not exist"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/instances/{name} [put]
func (h *Handler) UpdateInstance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		validatedName, err := validation.ValidateInstanceName(name)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_instance_name", err.Error())
			return
		}

