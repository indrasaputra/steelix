package steelix_test

import (
	"testing"
	"time"

	"github.com/indrasaputra/steelix"
)

type mockBackoff struct{}

func (mockBackoff) NextInterval() time.Duration { return 0 }

func TestNewClient(t *testing.T) {

}

func createRetryConfig(n uint32) *steelix.RetryConfig {
	return &steelix.RetryConfig{
		Backoff:  mockBackoff{},
		MaxRetry: n,
	}
}
