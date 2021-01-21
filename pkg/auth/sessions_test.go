package auth

import (
	"testing"
	"time"
)

func TestSessionStoreLifecycle(t *testing.T) {
	store := NewSessionStore(time.Hour)

	sess := store.Create("user-1", "Alice", "alice@example.com")
	if sess.ID == "" {
		t.Fatal("expected session ID to be set")
	}

	got := store.Get(sess.ID)
	if got == nil {
		t.Fatal("expected session to be found")
	}
