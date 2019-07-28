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

func TestGetFs(t *testing.T) {
	t.Run("nil server", func(t *testing.T) {
		var s *HTTPServer
		gfs := s.getFs()
		assert.Nil(t, gfs)
	})
	t.Run("non-nil server", func(t *testing.T) {
		mfs := fs.New()
		s := &HTTPServer{
			fs: mfs,
		}
		gfs := s.getFs()
		assert.Equal(t, mfs, gfs)
	})
}

func TestGetLayer(t *testing.T) {
	t.Run("nil server", func(t *testing.T) {
		var s *HTTPServer
		layer := s.getLayer()
		assert.Nil(t, layer)
	})
	t.Run("non-nil server", func(t *testing.T) {
		mlayer := new(mockLayer)
		s := &HTTPServer{
			layer: mlayer,
		}
		layer := s.getLayer()
		assert.Equal(t, mlayer, layer)
	})
}

func TestGetAddress(t *testing.T) {
	t.Run("nil server", func(t *testing.T) {
		var s *HTTPServer
		layer := s.getAddress()
		assert.Empty(t, layer)
	})
	t.Run("non-nil server", func(t *testing.T) {
		inaddr := "localhost"
		s := &HTTPServer{
			Address: inaddr,
		}
		outaddr := s.getAddress()
		assert.Equal(t, inaddr, outaddr)
	})
}

func TestServeHTTP(t *testing.T) {
	s := &HTTPServer{}
	testFunc := func() {
		s.ServeHTTP(nil, nil)
	}
	assert.NotPanics(t, testFunc)
}
