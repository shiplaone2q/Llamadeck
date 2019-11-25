package server

import (
	"context"
	"errors"
	"llamadeck/pkg/auth"
	"llamadeck/pkg/config"
	"llamadeck/pkg/database"
	"testing"
)

// stubAuthStore satisfies database.AuthStore for ResolveInstancePermission tests.
// Only GetInstancePermission is exercised here; the embedded nil interface makes the
// other methods available for interface satisfaction without implementation.
type stubAuthStore struct {
	database.AuthStore
	perm *auth.KeyPermission
	err  error
}

func (s *stubAuthStore) GetInstancePermission(ctx context.Context, keyID, instanceID int) (*auth.KeyPermission, error) {
	return s.perm, s.err
}

func newResolveMiddleware(store database.AuthStore) *APIAuthMiddleware {
	return NewAPIAuthMiddleware(config.AuthConfig{}, store)
