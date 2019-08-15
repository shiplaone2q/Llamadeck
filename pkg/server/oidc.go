package server

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"llamadeck/pkg/auth"
	"llamadeck/pkg/config"
)

const (
	sessionCookieName = "llamadeck_session"
	stateCookieName   = "llamadeck_oidc_state"
	callbackPath      = "/api/v1/auth/oidc/callback"
	stateCookieTTL    = 5 * time.Minute
)

// OIDCService implements browser login against an OpenID Connect provider.
// API keys remain the only way to authenticate non-browser clients.
type OIDCService struct {
	verifier *oidc.IDTokenVerifier
	oauth2   oauth2.Config
	cfg      config.OIDCConfig
	stateKey []byte

	Sessions *auth.SessionStore
}

// NewOIDCService creates a new OIDCService with eager IdP discovery.
func NewOIDCService(authCfg config.AuthConfig) (*OIDCService, error) {
	cfg := authCfg.OIDC
	if !cfg.Enabled() {
		return nil, fmt.Errorf("oidc is not configured: issuer_url and client_id are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("discovering OIDC provider at %s: %w "+
			"(check that the issuer URL is reachable and serves .well-known/openid-configuration)", cfg.IssuerURL, err)
	}

	stateKey := make([]byte, 32)
	if _, err := rand.Read(stateKey); err != nil {
		return nil, fmt.Errorf("generating OIDC state key: %w", err)
	}

	return &OIDCService{
		verifier: provider.Verifier(&oidc.Config{ClientID: cfg.ClientID}),
		oauth2: oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint:     provider.Endpoint(),
			Scopes:       cfg.Scopes,
		},
		cfg:      cfg,
		stateKey: stateKey,
		Sessions: auth.NewSessionStore(cfg.SessionTTL),
	}, nil
}

// redirectURL resolves the OAuth2 callback URL for this request, preferring
// the configured override and otherwise deriving it from forwarded headers.
func (s *OIDCService) redirectURL(r *http.Request) string {
	if s.cfg.RedirectURL != "" {
		return s.cfg.RedirectURL
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}
	host := r.Host
	if forwarded := r.Header.Get("X-Forwarded-Host"); forwarded != "" {
		host = forwarded
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, callbackPath)
}

// stateCookiePath returns the Path attribute for the state cookie: the
// browser-facing path of the callback endpoint. Behind a subpath proxy the
// browser reaches the callback under the external base prefix, so the cookie
// must be scoped to that path or the browser will not send it back.
func (s *OIDCService) stateCookiePath(r *http.Request) string {
	u, err := url.Parse(s.redirectURL(r))
	if err != nil || u.Path == "" {
		return callbackPath
	}
	return u.Path
}

// webRootPath returns the path the callback should send the browser to
func (s *OIDCService) webRootPath(r *http.Request) string {
	u, err := url.Parse(s.redirectURL(r))
	if err != nil {
		return "/"
	}
	base, ok := strings.CutSuffix(u.Path, callbackPath)
	if !ok || base == "" {
		return "/"
