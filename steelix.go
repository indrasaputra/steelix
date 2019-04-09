// Package steelix wraps http client and makes it more resilient.
package steelix

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

const (
	serverErrCode     = 500
	maxPassedRequests = 5
)

var (
	errServer = errors.New("server replied with 5xx status code")
)

// Backoff is a contract for implementing backoff strategy.
type Backoff interface {
	// NextInterval returns the interval for the subsequent requests.
	NextInterval() time.Duration
}

// RetryConfig holds configuration for implementing retry strategy.
type RetryConfig struct {
	// Backoff is backoff strategy.
	Backoff Backoff
	// MaxRetry sets how many times a request should be tried if error happens.
	// To make it clear, this is how the request life cycle.
	// Say we set MaxRetry to 1.
	// First, a request will be launched. If an error occurred, then the request will be tried once again.
	// In other words, by setup MaxRetry to 1, at most there will be two trials.
	MaxRetry uint32
}

// BreakerConfig holds configuration for implementing circuit breaker strategy.
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

// Client wraps native golang http.Client
// but imbued by retry and circuit breaker strategy if supplied.
type Client struct {
	client        *http.Client
	retryConfig   *RetryConfig
	breakerConfig *BreakerConfig
	breaker       *gobreaker.CircuitBreaker
}

// NewClient creates an instance of steelix.Client.
// If RetryConfig is not set, the default retry config will be applied.
// If BreakerConfig is not set, the default breaker config will be applied.
//
// Default retry config is:
//   - Backoff: NoBackoff
//   - MaxRetry: 0
//
// Default breaker config is:
//   - Name: "steelix-client"
//   - MinRequests: 10
//   - MinConsecutiveFailures: 10
//   - FailurePercentage: 25
func NewClient(client *http.Client, rc *RetryConfig, bc *BreakerConfig) *Client {
	rc = buildRetryConfig(rc)
}

// Do does almost the same things http.Client.Do does.
// The differences are resiliency strategies.
// While the native http.Client.Do only sends a request and returns a response,
// this method wraps it with resiliency strategies,
// such as retry and circuit breaker.
//
// For example, when RetryConfig is set, the failed request will be repeated until max retry is exceeded.
// Before sending a request, a header X-Steelix-Retry will be set to the request.
// Its value is the current retry count.
//
// When BreakerConfig is set, the request will be launched inside circuit breaker.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
		tmp  interface{}
	)

	for i := uint32(0); i <= c.retryConfig.MaxRetry; i++ {
		if resp != nil {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}

		req.Header.Set("X-Steelix-Retry", fmt.Sprintf("%d", i))

		if c.breaker == nil {
			resp, err = c.client.Do(req)
		} else {
			tmp, err = c.breaker.Execute(func() (interface{}, error) {
				r, e := c.client.Do(req)
				if r != nil && r.StatusCode >= serverErrCode {
					return r, errServer
				}
				return r, e
			})
			if tmp != nil {
				resp = tmp.(*http.Response)
			}
		}

		if err != nil || resp.StatusCode >= serverErrCode {
			time.Sleep(c.retryConfig.Backoff.NextInterval())
			continue
		}
		break
	}

	return resp, err
}

func buildRetryConfig(rc *RetryConfig) *RetryConfig {
	if rc == nil {
		rc = &RetryConfig{
			Backoff:  &noBackoff{},
			MaxRetry: 0,
		}
	}
	return rc
}
