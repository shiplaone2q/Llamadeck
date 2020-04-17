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
