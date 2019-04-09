package steelix

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

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
	var (
		resp *http.Response
		err  error
		tmp  interface{}
	)

	for i := uint32(0); i <= h.client.config.MaxRetry; i++ {
		if resp != nil {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}

		req.Header.Set("X-Steelix-Retry", fmt.Sprintf("%d", i))
		tmp, err = h.breaker.Execute(func() (interface{}, error) {
			r, e := h.client.client.Do(req)
			if r != nil && r.StatusCode >= 500 {
				return r, errServer
			}
			return r, e
		})
		if tmp != nil {
			resp = tmp.(*http.Response)
		}
		if err != nil {
			time.Sleep(h.client.config.Backoff.NextInterval())
			continue
		}
		break
	}

	return resp, err
}
