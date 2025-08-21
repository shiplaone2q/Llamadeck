package server

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"llamadeck/pkg/config"
)

// fakeIdP is a minimal OpenID Connect provider: discovery, a JWKS endpoint,
// and a token endpoint that mints RS256 id_tokens. The nonce embedded in
// issued tokens is test-controlled so the callback's nonce check can be
// exercised both ways.
type fakeIdP struct {
	server   *httptest.Server
	key      *rsa.PrivateKey
	clientID string
	nonce    string
}

func newFakeIdP(t *testing.T) *fakeIdP {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generating signing key: %v", err)
	}

	idp := &fakeIdP{key: key, clientID: "llamadeck-test"}
	mux := http.NewServeMux()

	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
		writeJSONTest(w, map[string]any{
			"issuer":                                idp.server.URL,
			"authorization_endpoint":                idp.server.URL + "/auth",
			"token_endpoint":                        idp.server.URL + "/token",
			"jwks_uri":                              idp.server.URL + "/keys",
			"response_types_supported":              []string{"code"},
			"subject_types_supported":               []string{"public"},
			"id_token_signing_alg_values_supported": []string{"RS256"},
		})
	})

	mux.HandleFunc("/keys", func(w http.ResponseWriter, _ *http.Request) {
		writeJSONTest(w, map[string]any{
			"keys": []map[string]any{{
				"kty": "RSA", "use": "sig", "alg": "RS256", "kid": "test-key",
				"n": base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
				"e": "AQAB",
			}},
		})
	})

	mux.HandleFunc("/token", func(w http.ResponseWriter, _ *http.Request) {
		token := signRS256Test(t, key, map[string]any{
			"iss":    idp.server.URL,
			"aud":    idp.clientID,
			"sub":    "user-1",
			"name":   "Alice",
			"email":  "alice@example.com",
			"groups": []string{"llamadeck-admins"},
			"exp":    time.Now().Add(time.Hour).Unix(),
			"iat":    time.Now().Unix(),
			"nonce":  idp.nonce,
		})
		writeJSONTest(w, map[string]any{
			"access_token": "fake-access-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
			"id_token":     token,
		})
	})

	idp.server = httptest.NewServer(mux)
	t.Cleanup(idp.server.Close)
