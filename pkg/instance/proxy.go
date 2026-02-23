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
		return nil, fmt.Errorf("instance %s has no options set", instance.Name)
	}

	if instance.IsRemote() {

		// Take the first remote node as the target for now
		var nodeName string
		for node := range options.Nodes {
			nodeName = node
			break
		}

		if nodeName == "" {
			return nil, fmt.Errorf("instance %s has no remote nodes defined", p.instance.Name)
		}

		node, ok := p.instance.globalNodesConfig[nodeName]
		if !ok {
			return nil, fmt.Errorf("remote node %s is not defined", nodeName)
		}

		p.targetURL, err = url.Parse(node.Address)
		if err != nil {
			return nil, fmt.Errorf("failed to parse target URL for remote instance %s: %w", p.instance.Name, err)
		}

		p.apiKey = node.APIKey
	} else {
		// Get host/port from process
