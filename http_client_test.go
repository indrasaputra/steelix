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

func createFailHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`not ok`))
	}
}
