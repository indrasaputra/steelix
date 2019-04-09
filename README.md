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