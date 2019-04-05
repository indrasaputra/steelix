package steelix_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/indrasaputra/steelix"
	"github.com/stretchr/testify/assert"
)

type mockBackoff struct{}

func (m mockBackoff) NextInterval() time.Duration {
	return 0
}

func TestNewHTTPClient(t *testing.T) {
	config := &steelix.ClientConfig{
		Backoff:  mockBackoff{},
		MaxRetry: 0,
	}
	client := steelix.NewHTTPClient(http.DefaultClient, config)

	assert.NotNil(t, client)
	assert.IsType(t, &steelix.HTTPClient{}, client)
}