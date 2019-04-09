package steelix_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/steelix"
)

func TestNewHTTPBreakerClient(t *testing.T) {
	sc := createSteelixHTTPClient()
	cfg := createBreakerConfig()

	client := steelix.NewHTTPBreakerClient(sc, cfg)
	assert.NotNil(t, client)
	assert.IsType(t, &steelix.HTTPBreakerClient{}, client)
}

func createSteelixHTTPClient() *steelix.HTTPClient {
	config := createConfig(1)
	return steelix.NewHTTPClient(http.DefaultClient, config)
}

func createBreakerConfig() *steelix.BreakerConfig {
	return &steelix.BreakerConfig{
		Name:                   "steelix breaker",
		MinRequests:            1,
		MinConsecutiveFailures: 2,
		FailurePercentage:      50,
	}
}
