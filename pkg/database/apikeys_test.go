package database

import (
	"context"
	"llamadeck/pkg/auth"
	"testing"
	"time"
)

// Keys attributed to SSO users (user_id != "system") must appear in the key
// listing; it is an admin view over all keys, not one user's keys.
func TestListKeysIncludesAllOwners(t *testing.T) {
