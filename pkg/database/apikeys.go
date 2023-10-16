package database

import (
	"context"
	"database/sql"
	"fmt"
	"llamadeck/pkg/auth"
	"time"
)

// CreateKey inserts a new API key with permissions (transactional)
func (db *sqliteDB) CreateKey(ctx context.Context, key *auth.APIKey, permissions []auth.KeyPermission) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert the API key
	query := `
		INSERT INTO api_keys (key_hash, name, user_id, permission_mode, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	var expiresAt sql.NullInt64
	if key.ExpiresAt != nil {
		expiresAt = sql.NullInt64{Int64: *key.ExpiresAt, Valid: true}
	}

	result, err := tx.ExecContext(ctx, query,
		key.KeyHash, key.Name, key.UserID, key.PermissionMode,
		expiresAt, key.CreatedAt, key.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert API key: %w", err)
	}

	keyID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	key.ID = int(keyID)

	// Insert permissions if per-instance mode
	if key.PermissionMode == auth.PermissionModePerInstance {
		for _, perm := range permissions {
			query := `
				INSERT INTO key_permissions (key_id, instance_id, can_start, can_evict)
				VALUES (?, ?, ?, ?)
			`
			_, err := tx.ExecContext(ctx, query, key.ID, perm.InstanceID, perm.CanStart, perm.CanEvict)
			if err != nil {
				return fmt.Errorf("failed to insert permission for instance %d: %w", perm.InstanceID, err)
			}
		}
	}

	return tx.Commit()
}

// GetKeyByID retrieves an API key by ID
func (db *sqliteDB) GetKeyByID(ctx context.Context, id int) (*auth.APIKey, error) {
	query := `
		SELECT id, key_hash, name, user_id, permission_mode, expires_at, created_at, updated_at, last_used_at
		FROM api_keys
		WHERE id = ?
	`

	var key auth.APIKey
	var expiresAt sql.NullInt64
	var lastUsedAt sql.NullInt64

	err := db.QueryRowContext(ctx, query, id).Scan(
		&key.ID, &key.KeyHash, &key.Name, &key.UserID, &key.PermissionMode,
		&expiresAt, &key.CreatedAt, &key.UpdatedAt, &lastUsedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("API key not found")
		}
		return nil, fmt.Errorf("failed to query API key: %w", err)
	}

	if expiresAt.Valid {
		key.ExpiresAt = &expiresAt.Int64
	}
	if lastUsedAt.Valid {
		key.LastUsedAt = &lastUsedAt.Int64
	}

	return &key, nil
}
