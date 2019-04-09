// Package steelix wraps http client and makes it more resilient.
package steelix

import "net/http"

// RetryConfig holds configuration for implementing retry strategy.
type RetryConfig struct {
	// Backoff is backoff strategy.
	Backoff Backoff
	// MaxRetry sets how many times a request should be tried if error happens.
	// To make it clear, this is how the request life cycle.
	// Say we set MaxRetry to 1.
	// First, a request will be launched. If an error occurred, then the request will be tried once again.
	// In other words, by setup MaxRetry to 1, at most there will be two trials.
	MaxRetry uint32
}

// Client wraps native golang http.Client
// but imbued by retry and circuit breaker strategy if supplied.
type Client struct {
	client *http.Client
}
