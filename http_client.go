// Package steelix wraps http client and makes it more resilient.
package steelix

import "time"

// Backoff is a contract for implementing backoff strategy.
type Backoff interface {
	// NextInterval returns the interval for the subsequent requests.
	NextInterval() time.Duration
}

// ClientConfig holds any configuration needed by HTTPClient.
type ClientConfig struct {
}

// HTTPClient wraps native golang http client.
// In addition, it provides retry strategy.
type HTTPClient struct {
}
