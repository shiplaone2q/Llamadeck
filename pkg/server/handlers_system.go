package server

import (
	"fmt"
	"net/http"
)

// VersionHandler godoc
// @Summary Get llamadeck version
// @Description Returns the version of the llamadeck command
// @Tags System
// @Security ApiKeyAuth
// @Produces text/plain
// @Success 200 {string} string "Version information"
// @Failure 500 {string} string "Internal Server Error"
