// Package middlewares provides common middleware handlers.
package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mateuszdyminski/go-profiling/stats"
	"github.com/varstr/uaparser"
)

// defaultVersion of REST api. Used when clients has no version in http header.
const defaultVersion = 3

// acceptVersions list supported REST client versions
var acceptVersions = []string{
	"vnd.app.v1",
	"vnd.app.v2",
	"vnd.app.v3",
}

func WithStats(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		tags := getStatsTags(r)

		stats.IncCounter("handler.received", tags, 1)

		responseStatistics := newResponseStats(w)

		// run endpoint here
		h(responseStatistics, r)

		stats.IncCounter(fmt.Sprintf("handler.httpCode.%d", responseStatistics.httpCode), tags, 1)

		stats.IncCounter(fmt.Sprintf("handler.apiVersion.%d", clientVersion(r)), tags, 1)

		duration := time.Since(start)
		stats.RecordTimer("handler.latency", tags, duration)
	}
}

func getStatsTags(r *http.Request) map[string]string {
	userBrowser, userOS := parseUserAgent(r.UserAgent())
	stats := map[string]string{
		"browser":  userBrowser,
		"os":       userOS,
		"endpoint": filepath.Base(r.URL.Path),
	}
	host, err := os.Hostname()
	if err == nil {
		if idx := strings.IndexByte(host, '.'); idx > 0 {
			host = host[:idx]
		}
		stats["host"] = host
	}
	return stats
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
// sample header application/vnd.app.v2+json
func clientVersion(req *http.Request) int {
	a := strings.Split(req.Header.Get("Accept"), "/")
	var v string
	if len(a) > 1 {
		versionAndType := strings.Split(a[1], "+")

		v = versionAndType[0]
		if len(versionAndType) > 1 {
			req.Header.Set("X-Version", strconv.Itoa(version(v)))
			req.Header.Set("X-Content", versionAndType[1])
		}
	}

	return version(v)
}

func version(v string) int {
	found := false
	for _, ver := range acceptVersions {
		if ver == v {
			found = true
			break
		}
	}

	if !found {
		return defaultVersion
	}

	ver, err := strconv.Atoi(v[9:])
	if err != nil {
		return defaultVersion
	}

	return ver
}