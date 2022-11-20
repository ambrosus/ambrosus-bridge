package rate_limiter

import (
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/time/rate"
)

const (
	DefaultLimitPeriod  = 1 * time.Second
	DefaultRequestCount = 5
)

// DialHTTP creates a new RPC client that connects to an RPC server over HTTP with default rate limiter options.
func DialHTTP(endpoint string) (*rpc.Client, error) {
	httpClient := new(http.Client)
	httpClient.Transport = NewThrottledTransport(DefaultLimitPeriod, DefaultRequestCount, http.DefaultTransport)
	return rpc.DialHTTPWithClient(endpoint, httpClient)
}

// DialHTTPWithLimitOptions creates a new RPC client that connects to an RPC server over HTTP with specified rate limiter options.
func DialHTTPWithLimitOptions(endpoint string, limitPeriod time.Duration, requestCount int) (*rpc.Client, error) {
	httpClient := new(http.Client)
	httpClient.Transport = NewThrottledTransport(limitPeriod, requestCount, http.DefaultTransport)
	return rpc.DialHTTPWithClient(endpoint, httpClient)
}

// Thanks @zdebra for https://gist.github.com/zdebra/10f0e284c4672e99f0cb767298f20c11

// ThrottledTransport Rate Limited HTTP Client
type ThrottledTransport struct {
	roundTripperWrap http.RoundTripper
	ratelimiter      *rate.Limiter
}

func (c *ThrottledTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	err := c.ratelimiter.Wait(r.Context()) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}
	return c.roundTripperWrap.RoundTrip(r)
}

// NewThrottledTransport wraps transportWrap with a rate limitter
// examle usage:
// client := http.DefaultClient
// client.Transport = NewThrottledTransport(10*time.Second, 60, http.DefaultTransport) allows 60 requests every 10 seconds
func NewThrottledTransport(limitPeriod time.Duration, requestCount int, transportWrap http.RoundTripper) http.RoundTripper {
	return &ThrottledTransport{
		roundTripperWrap: transportWrap,
		ratelimiter:      rate.NewLimiter(rate.Every(limitPeriod), requestCount),
	}
}
