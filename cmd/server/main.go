package main

import (
	"context"
	"fmt"
	"llamadeck/pkg/config"
	"llamadeck/pkg/database"
	"llamadeck/pkg/manager"
	"llamadeck/pkg/models"
	"llamadeck/pkg/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// version is set at build time using -ldflags "-X main.version=1.0.0"
var version string = "unknown"
var commitHash string = "unknown"
var buildTime string = "unknown"

// @title llamadeck API
// @version 1.0
// @description llamadeck is a control server for managing Llama Server instances.
// @license.name MIT License
// @license.url https://opensource.org/license/mit/
// @basePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {

	// --version flag to print the version
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("llamadeck version: %s\n", version)
		fmt.Printf("Commit hash: %s\n", commitHash)
		fmt.Printf("Build time: %s\n", buildTime)
		return
	}

	configPath := os.Getenv("LLAMADECK_CONFIG_PATH")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("Error loading config: %v\nUsing default configuration.", err)
	}

	// Set version information
	cfg.Version = version
	cfg.CommitHash = commitHash
	cfg.BuildTime = buildTime

	// Create data directory if it doesn't exist
	if cfg.Instances.AutoCreateDirs {
