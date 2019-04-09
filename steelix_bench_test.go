package steelix_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/indrasaputra/steelix"
)

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

func BenchmarkSteelixClient_Do_WithRetry(b *testing.B) {
	rc := createRetryConfig(5)
	client := steelix.NewClient(http.DefaultClient, rc, nil)
	benchmarkClient(b, client)
}

func BenchmarkSteelixClient_Do_WithBreaker(b *testing.B) {
	client := steelix.NewClient(http.DefaultClient, nil, createConsecutiveBreakerConfig())
	benchmarkClient(b, client)
}

func benchmarkClient(b *testing.B, client doer) {
	server := httptest.NewServer(http.HandlerFunc(createOkHandler()))
	defer server.Close()

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

	for i := 0; i < b.N; i++ {
		resp, _ := client.Do(req)
		if resp != nil {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}
	}
}
