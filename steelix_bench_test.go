package steelix_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/indrasaputra/steelix"
)

func BenchmarkClient_Do_WithRetry(b *testing.B) {
	rc := createRetryConfig(5)
	client := steelix.NewClient(http.DefaultClient, rc, nil)

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
