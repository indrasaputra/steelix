// Package steelix wraps http client and makes it more resilient.
package steelix

import "net/http"

// Client wraps native golang http.Client
// but imbued by retry and circuit breaker strategy if supplied.
type Client struct {
	client *http.Client
}
