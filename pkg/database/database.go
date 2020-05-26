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
