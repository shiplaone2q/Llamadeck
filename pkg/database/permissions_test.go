package database

import (
	"context"
	"llamadeck/pkg/auth"
	"testing"
	"time"
)

// newTestDB opens a fresh in-memory database with migrations applied. A single
// open connection guarantees every query shares the same in-memory database.
func newTestDB(t *testing.T) *sqliteDB {
	t.Helper()
	db, err := Open(&Config{Path: ":memory:", MaxOpenConnections: 1})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	return db
}

func insertTestInstance(t *testing.T, db *sqliteDB, id int, name string) {
	t.Helper()
	now := time.Now().Unix()
	_, err := db.ExecContext(context.Background(),
		`INSERT INTO instances (id, name, status, created_at, updated_at, options_json) VALUES (?, ?, 'stopped', ?, ?, '{}')`,
		id, name, now, now)
	if err != nil {
		t.Fatalf("insert instance %q: %v", name, err)
	}
}

func TestGetInstancePermission(t *testing.T) {
	db := newTestDB(t)
	insertTestInstance(t, db, 10, "model-a")

	now := time.Now().Unix()
	key := &auth.APIKey{
		KeyHash:        "hash",
		Name:           "k1",
		UserID:         "system",
		PermissionMode: auth.PermissionModePerInstance,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	perms := []auth.KeyPermission{
		{InstanceID: 10, CanStart: true, CanEvict: false},
	}
	if err := db.CreateKey(context.Background(), key, perms); err != nil {
		t.Fatalf("CreateKey: %v", err)
	}

