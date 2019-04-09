[![Build Status](https://travis-ci.com/indrasaputra/steelix.svg?branch=master)](https://travis-ci.com/indrasaputra/steelix)
[![codecov](https://codecov.io/gh/indrasaputra/steelix/branch/master/graph/badge.svg)](https://codecov.io/gh/indrasaputra/steelix)
[![GolangCI](https://golangci.com/badges/github.com/indrasaputra/steelix.svg)](https://golangci.com)
[![Go Report Card](https://goreportcard.com/badge/github.com/indrasaputra/steelix)](https://goreportcard.com/report/github.com/indrasaputra/steelix)
[![Documentation](https://godoc.org/github.com/indrasaputra/steelix?status.svg)](http://godoc.org/github.com/indrasaputra/steelix)

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
