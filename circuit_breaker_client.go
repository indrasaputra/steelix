package steelix

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

const (
	maxPassedRequests = 5
)

// BreakerConfig holds any configuration needed by HTTPBreakerConfig.
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

// HTTPBreakerClient wraps HTTPClient with circuit breaker functionality.
// It does what HTTPClient does and adds circuit breaker when doing its job.
type HTTPBreakerClient struct {
	client  *HTTPClient
	config  *BreakerConfig
	breaker *gobreaker.CircuitBreaker
}

// NewHTTPBreakerClient return an instance of HTTPBreakerClient.
func NewHTTPBreakerClient(client *HTTPClient, config *BreakerConfig) *HTTPBreakerClient {
	st := createBreakerSettings(config)
	breaker := gobreaker.NewCircuitBreaker(st)

	return &HTTPBreakerClient{
		client:  client,
		config:  config,
		breaker: breaker,
	}
}

// Do does almost the same thing as HTTPClient or Golang http.Client does.
// In addition, it adds circuit breaker functionality.
//
// When the ClientConfig is set, it also apply all resiliency strategies
// configured there, such as retry strategy.
func (h *HTTPBreakerClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := uint32(0); i <= h.client.config.MaxRetry; i++ {
		if resp != nil {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}

		req.Header.Set("X-Steelix-Retry", fmt.Sprintf("%d", i))
		tmp, err := h.breaker.Execute(func() (interface{}, error) {
			return h.client.client.Do(req)
		})
		if err != nil {
			time.Sleep(h.client.config.Backoff.NextInterval())
			continue
		}
		resp = tmp.(*http.Response)
		if resp.StatusCode >= 500 {
			time.Sleep(h.client.config.Backoff.NextInterval())
			continue
		}
		break
	}

	return resp, err
}

func createBreakerSettings(config *BreakerConfig) gobreaker.Settings {
	return gobreaker.Settings{
		Name:        config.Name,
		MaxRequests: maxPassedRequests,
		Interval:    0,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return readyToTrip(counts, config)
		},
	}
}

func readyToTrip(counts gobreaker.Counts, config *BreakerConfig) bool {
	if counts.Requests >= config.MinRequests && counts.ConsecutiveFailures >= config.MinConsecutiveFailures {
		return true
	}

	percentage := (float64(counts.TotalFailures) / float64(counts.Requests)) * 100
	if counts.Requests >= config.MinRequests && percentage >= config.FailurePercentage {
		return true
	}
	return false
}
