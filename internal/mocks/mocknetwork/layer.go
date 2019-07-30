package network

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockLayer struct {
	mock.Mock
}

func (m *MockLayer) ListenAndServe(addr string, handler http.Handler) error {
	args := m.Called(addr, handler)
	return args.Error(0)
}

func (m *MockLayer) ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	args := m.Called(addr, certFile, keyFile, handler)
	return args.Error(0)
}

func (m *MockLayer) FileServer(dir http.FileSystem) http.Handler {
	args := m.Called(dir)
	return args.Get(0).(http.Handler)
}
