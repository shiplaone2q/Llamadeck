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
	db := newTestDB(t)
	now := time.Now().Unix()

	for _, owner := range []string{"system", "alice@example.com"} {
		key := &auth.APIKey{
			KeyHash:        "hash-" + owner,
			Name:           "key-" + owner,
			UserID:         owner,
			PermissionMode: auth.PermissionModeAllowAll,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := db.CreateKey(context.Background(), key, nil); err != nil {
			t.Fatalf("creating key for %q: %v", owner, err)
		}
	}

	keys, err := db.ListKeys(context.Background())
	if err != nil {
