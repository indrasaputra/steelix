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
	config := createConfig(0)
	client := steelix.NewHTTPClient(http.DefaultClient, config)

	assert.NotNil(t, client)
	assert.IsType(t, &steelix.HTTPClient{}, client)
}

func TestHTTPClient_Do(t *testing.T) {

}

func createConfig(n uint32) *steelix.ClientConfig {
	return &steelix.ClientConfig{
		Backoff:  mockBackoff{},
		MaxRetry: n,
	}
}

func createOkHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`ok`))
	}
}
