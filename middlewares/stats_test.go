package middlewares

import (
	"net/http"
	"testing"
	"net/http/httptest"
	"fmt"
	"log"
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
	testCases := map[string]int {
		"application/vnd.app.v2+json" : 2,
		"application/vnd.app.v1" : 1,
		"application/wrong" : 3,
		"wrong" : 3,
		"" : 3,
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