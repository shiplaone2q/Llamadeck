package manager_test

import (
	"errors"
	"fmt"
	"llamadeck/pkg/backends"
	"llamadeck/pkg/config"
	"llamadeck/pkg/database"
	"llamadeck/pkg/instance"
	"llamadeck/pkg/manager"
	"strings"
	"testing"
	"time"
)

func TestCreateInstance_FailsWithDuplicateName(t *testing.T) {
	mngr := createTestManager(t)
	options := &instance.Options{
		BackendOptions: backends.Options{
			BackendType: backends.BackendTypeLlamaCpp,
			LlamaServerOptions: &backends.LlamaServerOptions{
				Model: "/path/to/model.gguf",
			},
		},
	}

	_, err := mngr.CreateInstance("test-instance", options)
	if err != nil {
		t.Fatalf("First CreateInstance failed: %v", err)
	}

	// Try to create duplicate
	_, err = mngr.CreateInstance("test-instance", options)
	if err == nil {
		t.Error("Expected error for duplicate instance name")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Expected duplicate name error, got: %v", err)
	}
}

func TestCreateInstance_FailsWhenMaxInstancesReached(t *testing.T) {
	appConfig := &config.AppConfig{
		Backends: config.BackendConfig{
			LlamaCpp: config.BackendSettings{
				Command: "llama-server",
			},
		},
		Instances: config.InstancesConfig{
			PortRange:            [2]int{8000, 9000},
			MaxInstances:         1, // Very low limit for testing
			TimeoutCheckInterval: 5,
		},
		Database: config.DatabaseConfig{
			Path:               ":memory:",
			MaxOpenConnections: 25,
			MaxIdleConnections: 5,
			ConnMaxLifetime:    5 * time.Minute,
		},
		LocalNode: "main",
		Nodes:     map[string]config.NodeConfig{},
	}
	db, err := database.Open(&database.Config{
		Path:               appConfig.Database.Path,
		MaxOpenConnections: appConfig.Database.MaxOpenConnections,
		MaxIdleConnections: appConfig.Database.MaxIdleConnections,
		ConnMaxLifetime:    appConfig.Database.ConnMaxLifetime,
	})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	limitedManager := manager.New(appConfig, db)

	options := &instance.Options{
		BackendOptions: backends.Options{
			BackendType: backends.BackendTypeLlamaCpp,
			LlamaServerOptions: &backends.LlamaServerOptions{
				Model: "/path/to/model.gguf",
			},
		},
	}

	_, err = limitedManager.CreateInstance("instance1", options)
	if err != nil {
		t.Fatalf("CreateInstance 1 failed: %v", err)
	}

	// This should fail due to max instances limit
	_, err = limitedManager.CreateInstance("instance2", options)
	if err == nil {
		t.Error("Expected error when exceeding max instances limit")
	}
	if !strings.Contains(err.Error(), "maximum number of instances") {
		t.Errorf("Expected max instances error, got: %v", err)
	}
}

func TestCreateInstance_FailsWithPortConflict(t *testing.T) {
	manager := createTestManager(t)

	options1 := &instance.Options{
		BackendOptions: backends.Options{
			BackendType: backends.BackendTypeLlamaCpp,
			LlamaServerOptions: &backends.LlamaServerOptions{
				Model: "/path/to/model.gguf",
				Port:  8080,
			},
		},
	}

	_, err := manager.CreateInstance("instance1", options1)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	// Try to create instance with same port
	options2 := &instance.Options{
		BackendOptions: backends.Options{
			BackendType: backends.BackendTypeLlamaCpp,
			LlamaServerOptions: &backends.LlamaServerOptions{
				Model: "/path/to/model2.gguf",
				Port:  8080, // Same port - should conflict
			},
		},
	}

	_, err = manager.CreateInstance("instance2", options2)
	if err == nil {
		t.Error("Expected error for port conflict")
	}
	if !strings.Contains(err.Error(), "port") && !strings.Contains(err.Error(), "in use") {
		t.Errorf("Expected port conflict error, got: %v", err)
	}
}

func TestInstanceOperations_FailWithNonExistentInstance(t *testing.T) {
	manager := createTestManager(t)

	options := &instance.Options{
		BackendOptions: backends.Options{
			BackendType: backends.BackendTypeLlamaCpp,
			LlamaServerOptions: &backends.LlamaServerOptions{
				Model: "/path/to/model.gguf",
			},
		},
	}

	_, err := manager.GetInstance("nonexistent")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
