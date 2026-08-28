package httpapi

import (
	"net/http"
	"sync/atomic"
)

type Metrics struct {
	Requests atomic.Int64
	Errors   atomic.Int64
}

func (m *Metrics) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Requests.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (m *Metrics) Snapshot() (int64, int64) { return m.Requests.Load(), m.Errors.Load() }
