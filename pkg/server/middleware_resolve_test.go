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
}

func withKey(ctx context.Context, key *auth.APIKey) context.Context {
	return context.WithValue(ctx, apiKeyContextKey, key)
}

func TestResolveInstancePermission(t *testing.T) {
	const instID = 7

	perInstanceKey := &auth.APIKey{ID: 1, PermissionMode: auth.PermissionModePerInstance}
	allowAllKey := &auth.APIKey{ID: 2, PermissionMode: auth.PermissionModeAllowAll}

	tests := []struct {
		name      string
		ctx       context.Context
		store     *stubAuthStore
		wantStart bool
		wantEvict bool
		wantErr   bool
	}{
		{
			name:      "management key (no APIKey in context) short-circuits to all-true",
			ctx:       context.Background(),
			store:     &stubAuthStore{},
			wantStart: true,
			wantEvict: true,
		},
		{
			name:      "allow_all key short-circuits to all-true",
			ctx:       withKey(context.Background(), allowAllKey),
			store:     &stubAuthStore{},
			wantStart: true,
			wantEvict: true,
		},
		{
			name:      "per_instance key both flags true",
			ctx:       withKey(context.Background(), perInstanceKey),
			store:     &stubAuthStore{perm: &auth.KeyPermission{KeyID: 1, InstanceID: instID, CanStart: true, CanEvict: true}},
			wantStart: true,
