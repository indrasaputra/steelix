[![Go Report Card](https://goreportcard.com/badge/github.com/indrasaputra/steelix)](https://goreportcard.com/report/github.com/indrasaputra/steelix)
[![Workflow](https://github.com/indrasaputra/steelix/workflows/Test/badge.svg)](https://github.com/indrasaputra/steelix/actions)
[![codecov](https://codecov.io/gh/indrasaputra/steelix/branch/main/graph/badge.svg?token=uLFqCaSQju)](https://codecov.io/gh/indrasaputra/steelix)
[![Maintainability](https://api.codeclimate.com/v1/badges/752f9d6d57202bf9b92e/maintainability)](https://codeclimate.com/github/indrasaputra/steelix/maintainability)
[![Go Reference](https://pkg.go.dev/badge/github.com/indrasaputra/steelix.svg)](https://pkg.go.dev/github.com/indrasaputra/steelix)
[![GolangCI](https://golangci.com/badges/github.com/indrasaputra/steelix.svg)](https://golangci.com)

# Steelix

Steelix is an HTTP client reinforcement using resiliency strategy.

## Description

Steelix wraps native golang HTTP client with some resiliency strategies. There are two resiliency strategies available, retry and circuit breaker.

## Installation

```
go get -u github.com/indrasaputra/steelix
```

## Usage

Struct `steelix.Client` wraps `http.Client`. Therefore, users should prepare their own `http.Client`, then use constructor to create an instance of `steelix.Client`.

To use retry and circuit breaker strategy, provide the respective configurations.

For more information, visit documentation in godoc.

```go
package main

import (
	"net/http"
	"time"

	"github.com/indrasaputra/backoff"
	"github.com/indrasaputra/steelix"
)

func main() {
	b := &backoff.ConstantBackoff{
		BackoffInterval: 200 * time.Millisecond,
		JitterInterval:  50 * time.Millisecond,
	}

	rc := &steelix.RetryConfig{
		Backoff:  b,
		MaxRetry: 3,
	}

	bc := &steelix.BreakerConfig{
		Name:                   "steelix-breaker",
		MinRequests:            10,
		MinConsecutiveFailures: 5,
		FailurePercentage:      20,
	}

	client := steelix.NewClient(http.DefaultClient, rc, bc)
	// omitted
}
```

Then, use `Do(req *http.Request)` method to send an HTTP request.

```go
req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
resp, err := client.Do(req)
```