// Package steelix wraps http client and makes it more resilient.
package steelix

import (
	"net/http"
	"time"
)

const (
	// DefaultMaxRetry is the default value for max retry.
	DefaultMaxRetry = 0
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
	MaxRetry int
}

// HTTPClient wraps native golang http client.
// In addition, it provides retry strategy.
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates an instance of HTTPClient.
func NewHTTPClient(client *http.Client, config *ClientConfig) *HTTPClient {
}
