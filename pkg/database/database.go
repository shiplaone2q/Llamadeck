package database

import (
	"context"
	"database/sql"
	"fmt"
	"llamadeck/pkg/auth"
	"llamadeck/pkg/instance"
	"log"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// InstanceStore defines interface for instance persistence operations
type InstanceStore interface {
	Save(inst *instance.Instance) error
	Delete(name string) error
	LoadAll() ([]*instance.Instance, error)
	Close() error
}

// AuthStore defines the interface for authentication operations
type AuthStore interface {
	CreateKey(ctx context.Context, key *auth.APIKey, permissions []auth.KeyPermission) error
	ListKeys(ctx context.Context) ([]*auth.APIKey, error)
	GetActiveKeys(ctx context.Context) ([]*auth.APIKey, error)
	GetKeyByID(ctx context.Context, id int) (*auth.APIKey, error)
	DeleteKey(ctx context.Context, id int) error
	TouchKey(ctx context.Context, id int) error
	GetPermissions(ctx context.Context, keyID int) ([]auth.KeyPermission, error)
	GetInstancePermission(ctx context.Context, keyID, instanceID int) (*auth.KeyPermission, error)
}

// Config contains database configuration settings
type Config struct {
	// Database file path (relative to data_dir or absolute)
	Path string

	// Connection settings
	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    time.Duration
}

// sqliteDB wraps database connection with configuration
type sqliteDB struct {
	*sql.DB
	config *Config
}

// Open creates a new database connection with provided configuration
func Open(config *Config) (*sqliteDB, error) {
	if config == nil {
		return nil, fmt.Errorf("database config cannot be nil")
	}

