package network

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockLayer is a mock of network.Layer
type MockLayer struct {
	mock.Mock
}

// ListenAndServe listens and serves HTTP
func (m *MockLayer) ListenAndServe(addr string, handler http.Handler) error {
	args := m.Called(addr, handler)
	return args.Error(0)
}

// ListenAndServeTLS listens and serves HTTPS
func (m *MockLayer) ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	args := m.Called(addr, certFile, keyFile, handler)
	return args.Error(0)
}

// FileServer serves a file system as a webserver
func (m *MockLayer) FileServer(dir http.FileSystem) http.Handler {
	args := m.Called(dir)
	return args.Get(0).(http.Handler)
}
