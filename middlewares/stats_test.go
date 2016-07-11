package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStats(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	WithStats(handler).ServeHTTP(w, req)

	fmt.Printf("%d - %s", w.Code, w.Body.String())
}

func TestClientVersion(t *testing.T) {
	testCases := map[string]int{
		"application/vnd.app.v2+json": 3,
		"application/vnd.app.v1":      1,
		"application/wrong":           3,
		"wrong":                       3,
		"":                            3,
	}

	for k, v := range testCases {
		req, err := http.NewRequest("GET", "http://example.com/foo", nil)
		if err != nil {
			t.Errorf("there should be no err here", err)
		}

		req.Header.Set("Accept", k)

		version := clientVersion(req)
		if version != v {
			t.Errorf("wrong client version. should be: %d, is: %d", v, version)
		}
	}
}

func BenchmarkClientVersion(b *testing.B) {
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Set("Accept", "application/vnd.app.v6")

	for i := 0; i < b.N; i++ {
		clientVersion(req)
	}
}
