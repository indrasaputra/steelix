// Package steelix wraps http client and makes it more resilient.
package steelix

import (
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

// Backoff is a contract for implementing backoff strategy.
type Backoff interface {
	// NextInterval returns the interval for the subsequent requests.
	NextInterval() time.Duration
}

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

// BreakerConfig holds configuration for implementing circuit breaker strategy.
type BreakerConfig struct {
	// Name is the name of circuit breaker.
	Name string
	// MinRequests is the minimum requests needed for breaker to start applying
	// the logic whether it needs to change state.
	//
	// If we set MinRequests = 10, the breaker will apply the logic
	// if there is at least 10 requests. Otherwise, the logic doesn't apply.
	MinRequests uint32
	// MinConsecutiveFailures is the minimum number of failed requests that will
	// make the breaker changes its states from closed to open or half-open.
	//
	// This configuration is used together with MinRequests, which means
	// if we set MinConsecutiveFailures=5 and MinRequests=10, then there are
	// 7 failed requests, the breaker will not change its state.
	//
	// The breaker will change from closed to open
	// if either MinConsecutiveFailures or FailurePercentage condition are met.
	MinConsecutiveFailures uint32
	// FailurePercentage is a percentage which will change the breaker state
	// from closed to open if the percentage of failure requests is equal or higher
	// than the given value.
	//
	// This configuration will run together with MinRequests and alongside
	// the MinConsecutiveFailures.
	//
	// The breaker will change from closed to open
	// if either MinConsecutiveFailures or FailurePercentage condition are met.
	FailurePercentage float64
}

// Client wraps native golang http.Client
// but imbued by retry and circuit breaker strategy if supplied.
type Client struct {
	client        *http.Client
	retryConfig   *RetryConfig
	breakerConfig *BreakerConfig
	breaker       *gobreaker.CircuitBreaker
}
