package network

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/mocklog"
	"github.com/LLKennedy/webserver/internal/mocks/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/tools/godoc/vfs"
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

func (m *mockLayer) FileServer(dir http.FileSystem) http.Handler {
	args := m.Called(dir)
	return args.Get(0).(http.Handler)
}

func TestNewHTTPServer(t *testing.T) {
	mfs := vfs.NewNameSpace()
	layer := HTTP{}
	logger := mocklog.New()
	s := NewHTTPServer(logger, mfs, layer)
	assert.Equal(t, &HTTPServer{logger: logger, Address: "localhost", layer: layer, fileServer: http.FileServer(network.NewDir(mfs))}, s)
	assert.Equal(t, "", logger.GetContents())
}

func TestStart(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mfs := vfs.NewNameSpace()
		layer := new(mockLayer)
		logger := mocklog.New()
		s := &HTTPServer{
			logger:     logger,
			layer:      layer,
			fileServer: http.FileServer(network.NewDir(mfs)),
		}
		layer.On("ListenAndServe", s.getAddress(), s).Return(nil)
		err := s.Start()
		assert.NoError(t, err)
		assert.Equal(t, "<nil>\n", logger.GetContents())
	})
	t.Run("error", func(t *testing.T) {
		mfs := vfs.NewNameSpace()
		layer := new(mockLayer)
		logger := mocklog.New()
		s := &HTTPServer{
			logger:     logger,
			layer:      layer,
			fileServer: http.FileServer(network.NewDir(mfs)),
		}
		layer.On("ListenAndServe", s.getAddress(), s).Return(fmt.Errorf("some network error"))
		err := s.Start()
		assert.EqualError(t, err, "http server closed unexpectedly: some network error")
		assert.Equal(t, "http server closed unexpectedly: some network error\n", logger.GetContents())
	})
}

func TestGetFs(t *testing.T) {
	t.Run("nil server", func(t *testing.T) {
		var s *HTTPServer
		gfs := s.getFs()
		assert.Nil(t, gfs)
	})
	t.Run("non-nil server", func(t *testing.T) {
		mfs := new(mockHandler)
		s := &HTTPServer{
			fileServer: mfs,
		}
		gfs := s.getFs()
		assert.Equal(t, mfs, gfs)
	})
}

func TestGetLogger(t *testing.T) {
	t.Run("nil server", func(t *testing.T) {
		var s *HTTPServer
		logger := s.getLogger()
		assert.Nil(t, logger)
	})
	t.Run("non-nil server", func(t *testing.T) {
		mlogger := log.New(os.Stdout, "test", log.Flags())
		s := &HTTPServer{
			logger: mlogger,
		}
		logger := s.getLogger()
		assert.Equal(t, mlogger, logger)
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

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {

}

func TestServeHTTP(t *testing.T) {
	mfs := vfs.NewNameSpace()
	s := &HTTPServer{
		fileServer: http.FileServer(network.NewDir(mfs)),
	}
	testFunc := func() {
		defer func() {
			if r := recover(); r != nil {
				assert.Failf(t, "caught error", "%v\n%s", r, debug.Stack())
			}
		}()
		s.ServeHTTP(new(mockResponseWriter), &http.Request{URL: &url.URL{}})
	}
	assert.NotPanics(t, testFunc)
}
