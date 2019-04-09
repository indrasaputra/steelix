package steelix_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
	config := createConfig(1)
	client := steelix.NewHTTPClient(http.DefaultClient, config)

	// === test against server ===
	tables := []struct {
		handler func(http.ResponseWriter, *http.Request)
		status  int
		retry   string
	}{
		{createOkHandler(), http.StatusOK, "0"},
		{createFailHandler(), http.StatusInternalServerError, "1"},
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
			assert.Equal(t, table.retry, resp.Header.Get("X-Steelix-Retry"))
		})
	}

	// === test when request is not valid ===
	t.Run("invalid request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "inval!t", nil)
		assert.Nil(t, err)

		server := httptest.NewServer(http.HandlerFunc(createOkHandler()))
		defer server.Close()

		resp, err := client.Do(req)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
	})
}

func createConfig(n uint32) *steelix.ClientConfig {
	return &steelix.ClientConfig{
		Backoff:  mockBackoff{},
		MaxRetry: n,
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
