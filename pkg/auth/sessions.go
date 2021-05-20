package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// Session is an authenticated browser session created by OIDC login.
type Session struct {
	ID        string
	Sub       string
	Name      string
	Email     string
	ExpiresAt time.Time
}

// SessionStore keeps sessions in memory. Sessions do not survive restarts
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	ttl      time.Duration
}

// NewSessionStore creates a store where sessions live for ttl.
func NewSessionStore(ttl time.Duration) *SessionStore {
	return &SessionStore{
