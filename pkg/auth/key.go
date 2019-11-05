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
