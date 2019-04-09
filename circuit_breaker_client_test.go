package steelix_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func TestHTTPBreakerClient_Do(t *testing.T) {
	sc := createSteelixHTTPClient()
	cfg := createBreakerConfig()

	client := steelix.NewHTTPBreakerClient(sc, cfg)

	// === test against server ===
	tables := []struct {
		handler func(http.ResponseWriter, *http.Request)
		status  int
	}{
		{createOkHandler(), http.StatusOK},
		{createFailHandler(), http.StatusInternalServerError},
	}

	for _, table := range tables {
		t.Run(fmt.Sprintf("server return %d", table.status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(table.handler))
			defer server.Close()

			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			assert.Nil(t, err)

			resp, err := client.Do(req)
			defer func() {
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}()

			assert.Nil(t, err)
			assert.Equal(t, table.status, resp.StatusCode)
		})
	}

	// === test when request is not valid ===
	t.Run("invalid request", func(t *testing.T) {
		c := steelix.NewHTTPBreakerClient(sc, cfg)

		req, err := http.NewRequest(http.MethodGet, "inval!t", nil)
		assert.Nil(t, err)

		server := httptest.NewServer(http.HandlerFunc(createOkHandler()))
		defer server.Close()

		resp, err := c.Do(req)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
	})

	// === capture error on consecutive failure ===
	t.Run("capture consecutive failure", func(t *testing.T) {
		cc := createConfig(5)
		hc := steelix.NewHTTPClient(http.DefaultClient, cc)
		bc := &steelix.BreakerConfig{
			Name:                   "steelix breaker consecutive failure",
			MinRequests:            2,
			MinConsecutiveFailures: 2,
			FailurePercentage:      90,
		}

		client := steelix.NewHTTPBreakerClient(hc, bc)

		server := httptest.NewServer(http.HandlerFunc(createFailHandler()))
		defer server.Close()

		req, err := http.NewRequest(http.MethodGet, server.URL, nil)
		assert.Nil(t, err)

		_, err = client.Do(req)
		assert.NotNil(t, err)
	})
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
