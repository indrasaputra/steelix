package steelix_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/indrasaputra/steelix"
	"github.com/stretchr/testify/assert"
)

type mockBackoff struct{}

func (mockBackoff) NextInterval() time.Duration { return 0 }

func TestNewClient(t *testing.T) {
	// --- with retry, without breaker ---
	rc := createRetryConfig(5)
	client := steelix.NewClient(http.DefaultClient, rc, nil)
	assert.NotNil(t, client)

	// --- with retry and breaker ---
	client = steelix.NewClient(http.DefaultClient, rc, createConsecutiveBreakerConfig())
	assert.NotNil(t, client)
}

func createRetryConfig(n uint32) *steelix.RetryConfig {
	return &steelix.RetryConfig{
		Backoff:  mockBackoff{},
		MaxRetry: n,
	}
}

func createConsecutiveBreakerConfig() *steelix.BreakerConfig {
	return &steelix.BreakerConfig{
		Name:                   "steelix-consecutive-breaker",
		MinRequests:            2,
		MinConsecutiveFailures: 2,
		FailurePercentage:      10,
	}
}

func createPercentageBreakerConfig() *steelix.BreakerConfig {
	return &steelix.BreakerConfig{
		Name:                   "steelix-consecutive-breaker",
		MinRequests:            2,
		MinConsecutiveFailures: 3,
		FailurePercentage:      10,
	}
}

func createOkHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Steelix-Retry", r.Header.Get("X-Steelix-Retry"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`ok`))
	}
}

func createFailHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Steelix-Retry", r.Header.Get("X-Steelix-Retry"))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`not ok`))
	}
}
