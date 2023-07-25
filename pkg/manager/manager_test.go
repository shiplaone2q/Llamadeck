package manager_test

import (
	"fmt"
	"llamadeck/pkg/backends"
	"llamadeck/pkg/config"
	"llamadeck/pkg/database"
	"llamadeck/pkg/instance"
	"llamadeck/pkg/manager"
	"sync"
	"testing"
	"time"
)

func TestManager_PersistsAndLoadsInstances(t *testing.T) {
	tempDir := t.TempDir()
	appConfig := createTestAppConfig(tempDir)
	// Use file-based database for this test since we need to persist across connections
	appConfig.Database.Path = tempDir + "/test.db"

	// Create instance and check database was created
	db1, err := database.Open(&database.Config{
		Path:               appConfig.Database.Path,
		MaxOpenConnections: appConfig.Database.MaxOpenConnections,
		MaxIdleConnections: appConfig.Database.MaxIdleConnections,
		ConnMaxLifetime:    appConfig.Database.ConnMaxLifetime,
	})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	if err := database.RunMigrations(db1); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	manager1 := manager.New(appConfig, db1)
	options := &instance.Options{
		BackendOptions: backends.Options{
			BackendType: backends.BackendTypeLlamaCpp,
			LlamaServerOptions: &backends.LlamaServerOptions{
				Model: "/path/to/model.gguf",
				Port:  8080,
			},
		},
	}

	_, err = manager1.CreateInstance("test-instance", options)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	// Shutdown first manager to close database connection
	manager1.Shutdown()

	// Load instances from database
	db2, err := database.Open(&database.Config{
		Path:               appConfig.Database.Path,
		MaxOpenConnections: appConfig.Database.MaxOpenConnections,
		MaxIdleConnections: appConfig.Database.MaxIdleConnections,
		ConnMaxLifetime:    appConfig.Database.ConnMaxLifetime,
	})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	if err := database.RunMigrations(db2); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	manager2 := manager.New(appConfig, db2)
	instances, err := manager2.ListInstances()
	if err != nil {
		t.Fatalf("ListInstances failed: %v", err)
	}
	if len(instances) != 1 {
		t.Fatalf("Expected 1 loaded instance, got %d", len(instances))
	}
	if instances[0].Name != "test-instance" {
		t.Errorf("Expected loaded instance name 'test-instance', got %q", instances[0].Name)
	}

