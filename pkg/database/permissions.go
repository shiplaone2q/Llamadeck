package database

import (
	"context"
	"database/sql"
	"fmt"
	"llamadeck/pkg/auth"
)

// GetPermissions retrieves all permissions for a key
func (db *sqliteDB) GetPermissions(ctx context.Context, keyID int) ([]auth.KeyPermission, error) {
	query := `
		SELECT key_id, instance_id, can_start, can_evict
		FROM key_permissions
		WHERE key_id = ?
		ORDER BY instance_id
	`

	rows, err := db.QueryContext(ctx, query, keyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query key permissions: %w", err)
