// Package steelix wraps http client and makes it more resilient.
package steelix

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Backoff is a contract for implementing backoff strategy.
type Backoff interface {
	// NextInterval returns the interval for the subsequent requests.
	NextInterval() time.Duration
}

// ClientConfig holds any configuration needed by HTTPClient.
type ClientConfig struct {
	// Backoff is backoff strategy.
	Backoff Backoff
	// MaxRetry sets how many times a request should be tried if error happens.
	// To make it clear, this is how the request life cycle.
	// Say we set MaxRetry to 1.
	// First, a request will be launched. If an error occurred, then the request will be tried once.
	// In other words, by setup MaxRetry to 1, at most there will be two trials.
	MaxRetry uint32
}

// HTTPClient wraps native golang http client.
// In addition, it provides retry strategy.
type HTTPClient struct {
	client *http.Client
	config *ClientConfig
}

// NewHTTPClient creates an instance of HTTPClient.
func NewHTTPClient(client *http.Client, config *ClientConfig) *HTTPClient {
	return &HTTPClient{
		client: client,
		config: config,
	}
}

// Do does almost the same thing as http.Client.Do does.
// The differences are the resiliency strategies.
// While the native http.Client.Do only sends a request and returns the response,
// this method wraps it with resiliency strategies.
//
// For example, when MaxRetry is set, the failed request will be repeated until max retry is exceeded.
//
// Before sending a request, a context will be added to the request.
// The parent of the added context is taken from the request itself, so the original context won't go.
func (h *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := uint32(0); i <= h.config.MaxRetry; i++ {
		if resp != nil {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}

		req.Header.Set("X-Steelix-Retry", fmt.Sprintf("%d", i))
		resp, err = h.client.Do(req)
		if err != nil {
			time.Sleep(h.config.Backoff.NextInterval())
			continue
		}
		if resp.StatusCode >= 500 {
			time.Sleep(h.config.Backoff.NextInterval())
			continue
		}
		break
	}

	return resp, err
}
