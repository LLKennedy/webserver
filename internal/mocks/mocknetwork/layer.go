package mocknetwork

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// Layer is a mock of network.Layer
type Layer struct {
	mock.Mock
}

// ListenAndServe listens and serves HTTP
func (m *Layer) ListenAndServe(addr string, handler http.Handler) error {
	args := m.Called(addr, handler)
	return args.Error(0)
}

// ListenAndServeTLS listens and serves HTTPS
func (m *Layer) ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	args := m.Called(addr, certFile, keyFile, handler)
	return args.Error(0)
}

// FileServer serves a file system as a webserver
func (m *Layer) FileServer(dir http.FileSystem) http.Handler {
	args := m.Called(dir)
	var handler http.Handler
	if arg := args.Get(0); arg != nil {
		handler = arg.(http.Handler)
	}
	return handler
}
