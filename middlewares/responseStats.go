package middlewares

import "net/http"

// responseStats holds various statistics associated with HTTP request-response pair.
type callStats struct {
	w        http.ResponseWriter
	httpCode int
	resSize  int64
}

func newResponseStats(w http.ResponseWriter) *callStats {
	return &callStats{w: w}
}

func (r *callStats) Header() http.Header {
	return r.w.Header()
}

func (r *callStats) WriteHeader(code int) {
	r.w.WriteHeader(code)
	r.httpCode = code
}

func (r *callStats) Write(p []byte) (n int, err error) {
	if r.httpCode == 0 {
		r.httpCode = http.StatusOK
	}
	n, err = r.w.Write(p)
	r.resSize += int64(n)
	return
}

func (r *callStats) StatusCode() int {
	return r.httpCode
}

func (r *callStats) ResponseSize() int64 {
	return r.resSize
}
