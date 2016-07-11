// Package middlewares provides common middleware handlers.
package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mateuszdyminski/go-profiling/stats"
	"github.com/varstr/uaparser"
)

// defaultVersion of REST api. Used when clients has no version in http header.
const defaultVersion = 3

// acceptVersions list supported REST client versions
var acceptVersions = map[string]int{
	"vnd.app.v1": 1,
	"vnd.app.v2": 2,
	"vnd.app.v3": 3,
	"vnd.app.v4": 4,
	"vnd.app.v5": 5,
	"vnd.app.v6": 6,
}

func WithStats(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		tags := getStatsTags(r)

		stats.IncCounter("handler.received", tags, 1)

		responseStatistics := newResponseStats(w)

		// run endpoint here
		h(responseStatistics, r)

		// store statistics
		stats.IncCounter(fmt.Sprintf("handler.httpCode.%d", responseStatistics.httpCode), tags, 1)

		stats.IncCounter(fmt.Sprintf("handler.apiVersion.%d", clientVersion(r)), tags, 1)

		stats.RecordTimer("handler.latency", tags, time.Since(start))
	}
}

var hostname = getHostname()

func getStatsTags(r *http.Request) map[string]string {
	userBrowser, userOS := parseUserAgent(r.UserAgent())
	stats := map[string]string{
		"browser":  userBrowser,
		"os":       userOS,
		"endpoint": filepath.Base(r.URL.Path),
		"host":     hostname,
	}

	return stats
}

func getHostname() string {
	host, err := os.Hostname()
	if err == nil {
		if idx := strings.IndexByte(host, '.'); idx > 0 {
			host = host[:idx]
		}
	}

	return host
}

func parseUserAgent(uaString string) (browser, os string) {
	ua := uaparser.Parse(uaString)

	if ua.Browser != nil {
		browser = ua.Browser.Name
	}
	if ua.OS != nil {
		os = ua.OS.Name
	}

	return browser, os
}

// acceptedVersion checks Accept header in request to parse API version.
// sample header application/vnd.app.v2
func clientVersion(req *http.Request) int {
	a := strings.Split(req.Header.Get("Accept"), "/")
	if len(a) < 2 {
		return defaultVersion
	}

	val, ok := acceptVersions[a[1]]
	if !ok {
		return defaultVersion
	}

	return val
}
