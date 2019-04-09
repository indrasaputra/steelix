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

func TestClient_Do_WithRetry(t *testing.T) {
	rc := createRetryConfig(1)
	client := steelix.NewClient(http.DefaultClient, rc, nil)

	tables := []struct {
		handler func(w http.ResponseWriter, r *http.Request)
		status  int
		retry   string
	}{
		{createOkHandler(), 200, "0"},
		{createFailHandler(), 500, "1"},
	}

	for _, table := range tables {
		t.Run(fmt.Sprintf("Do with retry, server returns %d", table.status), func(t *testing.T) {
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
}

func TestClient_Do_WithRetryAndBreaker(t *testing.T) {

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
		MinRequests:            10,
		MinConsecutiveFailures: 10,
		FailurePercentage:      20,
	}
}

func createPercentageBreakerConfig() *steelix.BreakerConfig {
	return &steelix.BreakerConfig{
		Name:                   "steelix-consecutive-breaker",
		MinRequests:            10,
		MinConsecutiveFailures: 15,
		FailurePercentage:      20,
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
