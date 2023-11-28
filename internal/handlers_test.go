package internal_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/trail-l31/gazes-proxy/internal"
)

func BenchmarkProxyHandler(b *testing.B) {
	req, err := http.NewRequest("GET", "/?url=https://google.com", nil)

	if err != nil {
		b.Fatal(err)
	}

	// Simulate a request/response recorder for the handler
	rr := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		internal.ProxyHandler(rr, req)
	}
}
