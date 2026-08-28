package httpapi

import (
	"net/http"
	"time"
)

func NewServer(handler http.Handler) *http.Server {
	return &http.Server{Addr: ":8080", Handler: handler, ReadHeaderTimeout: 3 * time.Second}
}

func RegisterMux(handler http.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	return mux
}
