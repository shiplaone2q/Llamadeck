package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"llamadeck/pkg/instance"
	"log"
	"time"
)

// instanceRow represents a row in the instances table
type instanceRow struct {
	ID          int
	Name        string
	Status      string
	CreatedAt   int64
	UpdatedAt   int64
	OptionsJSON string
	OwnerUserID sql.NullString
}

// Create inserts a new instance into the database
func (db *sqliteDB) Create(ctx context.Context, inst *instance.Instance) error {
	if inst == nil {
		return fmt.Errorf("instance cannot be nil")
	}

	opts := inst.GetOptions()
	if opts == nil {
		return fmt.Errorf("instance options cannot be nil")
	}

	// Convert instance to database row
	row, err := db.instanceToRow(inst)
	if err != nil {
		return fmt.Errorf("failed to convert instance to row: %w", err)
	}

	// Insert into database
	query := `
		INSERT INTO instances (
			name, status, created_at, updated_at, options_json, owner_user_id
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := db.DB.ExecContext(ctx, query,
		row.Name, row.Status, row.CreatedAt, row.UpdatedAt, row.OptionsJSON, row.OwnerUserID,
	)

	if err != nil {
		return fmt.Errorf("failed to insert instance: %w", err)
	}

	// Get the auto-generated ID and set it on the instance
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	inst.ID = int(id)

	return nil
}

// GetByName retrieves an instance by name
func (db *sqliteDB) GetByName(ctx context.Context, name string) (*instance.Instance, error) {
	query := `
		SELECT id, name, status, created_at, updated_at, options_json, owner_user_id
		FROM instances
		WHERE name = ?
	`

	var row instanceRow
	err := db.DB.QueryRowContext(ctx, query, name).Scan(
		&row.ID, &row.Name, &row.Status, &row.CreatedAt, &row.UpdatedAt, &row.OptionsJSON, &row.OwnerUserID,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("instance not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query instance: %w", err)
	}

	return db.rowToInstance(&row)
}

// GetAll retrieves all instances from the database
func (db *sqliteDB) GetAll(ctx context.Context) ([]*instance.Instance, error) {
	query := `
		SELECT id, name, status, created_at, updated_at, options_json, owner_user_id
		FROM instances
		ORDER BY created_at ASC
	`

	rows, err := db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query instances: %w", err)
	}
	defer rows.Close()

	var instances []*instance.Instance
	for rows.Next() {
		var row instanceRow
		err := rows.Scan(
			&row.ID, &row.Name, &row.Status, &row.CreatedAt, &row.UpdatedAt, &row.OptionsJSON, &row.OwnerUserID,
		)
		if err != nil {
			log.Printf("Failed to scan instance row: %v", err)
			continue
		}

		inst, err := db.rowToInstance(&row)
		if err != nil {
			log.Printf("Failed to convert row to instance: %v", err)
			continue
		}

		instances = append(instances, inst)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return instances, nil
}

// Update updates an existing instance
func (db *sqliteDB) Update(ctx context.Context, inst *instance.Instance) error {
	if inst == nil {
		return fmt.Errorf("instance cannot be nil")
	}

	opts := inst.GetOptions()
	if opts == nil {
		return fmt.Errorf("instance options cannot be nil")
	}

	// Convert instance to database row
	row, err := db.instanceToRow(inst)
	if err != nil {
