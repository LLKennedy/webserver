package network

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLayer struct {
	mock.Mock
}

func (m *mockLayer) ListenAndServe(addr string, handler http.Handler) error {
	args := m.Called(addr, handler)
	return args.Error(0)
}

func (m *mockLayer) ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	args := m.Called(addr, certFile, keyFile, handler)
	return args.Error(0)
}

func TestNewHTTPServer(t *testing.T) {
	mfs := fs.New()
	layer := HTTP{}
	s := NewHTTPServer(mfs, layer)
	assert.Equal(t, &HTTPServer{Address: "localhost", fs: mfs, layer: layer}, s)
}

func TestStart(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mfs := fs.New()
		layer := new(mockLayer)
		s := &HTTPServer{
			layer: layer,
			fs:    mfs,
		}
		layer.On("ListenAndServe", s.getAddress(), s).Return(nil)
		err := s.Start()
		assert.NoError(t, err)
	})
	t.Run("error", func(t *testing.T) {
		mfs := fs.New()
		layer := new(mockLayer)
		s := &HTTPServer{
			layer: layer,
			fs:    mfs,
		}
		layer.On("ListenAndServe", s.getAddress(), s).Return(fmt.Errorf("some network error"))
		err := s.Start()
		assert.EqualError(t, err, "http server closed unexpectedly: some network error")
	})
}
