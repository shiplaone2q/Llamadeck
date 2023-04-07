package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type PermissionMode string

const (
	PermissionModeAllowAll    PermissionMode = "allow_all"
	PermissionModePerInstance PermissionMode = "per_instance"
)

type APIKey struct {
	ID             int
	KeyHash        string
	Name           string
	UserID         string
	PermissionMode PermissionMode
	ExpiresAt      *int64
	CreatedAt      int64
	UpdatedAt      int64
	LastUsedAt     *int64
}

type KeyPermission struct {
	KeyID      int
	InstanceID int
	CanStart   bool
	CanEvict   bool
}

// GenerateKey generates a cryptographically secure API key with the given prefix
func GenerateKey(prefix string) (string, error) {
