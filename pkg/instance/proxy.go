package instance

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// TimeProvider interface allows for testing with mock time
type TimeProvider interface {
	Now() time.Time
}

// realTimeProvider implements TimeProvider using the actual time
type realTimeProvider struct{}

func (realTimeProvider) Now() time.Time {
	return time.Now()
}

// proxy manages HTTP reverse proxy and request tracking for an instance.
type proxy struct {
	instance *Instance

	targetURL *url.URL
	apiKey    string // For remote instances

	responseHeaders map[string]string

	mu sync.RWMutex

	proxy     *httputil.ReverseProxy
	proxyOnce sync.Once
	proxyErr  error

	lastRequestTime  atomic.Int64
	inflightRequests atomic.Int32
	timeProvider     TimeProvider
}

// newProxy creates a new Proxy for the given instance
func newProxy(instance *Instance) (*proxy, error) {

	p := &proxy{
		instance:     instance,
		timeProvider: realTimeProvider{},
	}

	var err error

	options := instance.GetOptions()
	if options == nil {
