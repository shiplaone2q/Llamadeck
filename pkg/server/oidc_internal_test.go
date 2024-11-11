package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"

	"llamadeck/pkg/auth"
	"llamadeck/pkg/config"
)

func TestSignVerifyStateValue(t *testing.T) {
	key := []byte("test-key")

	signed := signStateValue(key, "some-state", "some-verifier", "some-nonce")

	verifier, nonce, ok := verifyStateValue(key, signed, "some-state")
	if !ok {
		t.Fatal("expected valid state cookie to verify")
	}
	if verifier != "some-verifier" {
		t.Errorf("verifier = %q, expected %q", verifier, "some-verifier")
	}
	if nonce != "some-nonce" {
		t.Errorf("nonce = %q, expected %q", nonce, "some-nonce")
	}
}

func TestVerifyStateValueRejectsBadInput(t *testing.T) {
	key := []byte("test-key")
	signed := signStateValue(key, "some-state", "some-verifier", "some-nonce")

	tests := []struct {
		name        string
		cookieValue string
		queryState  string
	}{
		{"state mismatch", signed, "other-state"},
		{"tampered signature", "some-state|some-verifier|some-nonce|deadbeef", "some-state"},
		{"garbage cookie", "garbage", "some-state"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, ok := verifyStateValue(key, tt.cookieValue, tt.queryState); ok {
				t.Error("expected verification to fail")
			}
		})
	}
}

func TestGroupsFromClaims(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		claim   string
		groups  []string
		present bool
		wantErr bool
	}{
		{"array of strings", `{"sub":"u1","groups":["admins","devs"]}`, "groups", []string{"admins", "devs"}, true, false},
		{"empty array", `{"groups":[]}`, "groups", []string{}, true, false},
		{"scalar string", `{"groups":"admins"}`, "groups", []string{"admins"}, true, false},
		{"claim absent", `{"sub":"u1","roles":["admins"]}`, "groups", nil, false, false},
		{"object shape", `{"groups":{"admins":true}}`, "groups", nil, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var raw map[string]json.RawMessage
			if err := json.Unmarshal([]byte(tt.payload), &raw); err != nil {
				t.Fatalf("fixture JSON invalid: %v", err)
			}

			groups, present, err := groupsFromClaims(raw, tt.claim)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error for a malformed claim")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if present != tt.present {
				t.Errorf("present = %v, expected %v", present, tt.present)
			}
			if !slices.Equal(groups, tt.groups) {
				t.Errorf("groups = %v, expected %v", groups, tt.groups)
