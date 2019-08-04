package network

import (
	"context"
	"fmt"
	"net/http"
)

// HTTP is the net/http implementation of the network.Layer interface
type HTTP struct {
	server http.Server
}

// Shutdown shuts down a server
func (h *HTTP) Shutdown(ctx context.Context) error {
	if h == nil {
		return fmt.Errorf("cannot shutdown nil server")
	}
	return h.server.Shutdown(ctx)
}

// ListenAndServe calls ListenAndServe on net/http
func (h *HTTP) ListenAndServe(addr string, handler http.Handler) error {
	if h == nil {
		return fmt.Errorf("cannot listen on nil server")
	}
	h.server.Addr = addr
	h.server.Handler = handler
	return h.server.ListenAndServe()
}

// ListenAndServeTLS calls ListenAndServeTLS on net/http
func (h *HTTP) ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	if h == nil {
		return fmt.Errorf("cannot listen on nil server")
	}
	h.server.Addr = addr
	h.server.Handler = handler
	return h.server.ListenAndServeTLS(certFile, keyFile)
}

// FileServer servies the file system over HTTP
func (h *HTTP) FileServer(fs http.FileSystem) http.Handler {
	return http.FileServer(fs)
}
