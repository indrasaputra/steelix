package steelix

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
	MinConsecutiveFailures uint32
}

// HTTPBreakerClient wraps HTTPClient with circuit breaker functionality.
// It does what HTTPClient does and adds circuit breaker when doing its job.
type HTTPBreakerClient struct {
	httpclient *HTTPClient
}
