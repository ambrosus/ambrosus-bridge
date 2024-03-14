package parity_rpc

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

// DialHTTP creates a new RPC client that connects to an RPC server over HTTP with default rate limiter options.
func DialHTTP(endpoint string) (*rpc.Client, error) {
	httpClient := new(http.Client)
	httpClient.Transport = &ReplaceInputToDataTransporter{roundTripperWrap: http.DefaultTransport}
	return rpc.DialHTTPWithClient(endpoint, httpClient)
}

// DialHTTPWithLimitOptions creates a new RPC client that connects to an RPC server over HTTP with specified rate limiter options.
func DialHTTPWithLimitOptions(endpoint string, limitPeriod time.Duration, requestCount int) (*rpc.Client, error) {
	httpClient := new(http.Client)
	httpClient.Transport = &ReplaceInputToDataTransporter{roundTripperWrap: http.DefaultTransport}
	return rpc.DialHTTPWithClient(endpoint, httpClient)
}

type ReplaceInputToDataTransporter struct {
	roundTripperWrap http.RoundTripper
}

func (c *ReplaceInputToDataTransporter) RoundTrip(r *http.Request) (*http.Response, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	bodyBytes = bytes.Replace(bodyBytes, []byte(`"input"`), []byte(`"data"`), -1)
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	r.ContentLength = int64(len(bodyBytes))
	return c.roundTripperWrap.RoundTrip(r)
}
