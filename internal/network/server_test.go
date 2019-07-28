package network

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/fs"
	"github.com/LLKennedy/webserver/internal/mocks/mocklog"
	"github.com/LLKennedy/webserver/internal/mocks/vnet"
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
	logger := mocklog.New()
	s := NewHTTPServer(logger, mfs, layer)
	assert.Equal(t, &HTTPServer{logger: logger, Address: "localhost", fs: mfs, layer: layer, hfs: vnet.NewDir(mfs)}, s)
	assert.Equal(t, "", logger.GetContents())
}

func TestStart(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mfs := fs.New()
		layer := new(mockLayer)
		logger := mocklog.New()
		s := &HTTPServer{
			logger: logger,
			layer:  layer,
			fs:     mfs,
		}
		layer.On("ListenAndServe", s.getAddress(), s).Return(nil)
		err := s.Start()
		assert.NoError(t, err)
		assert.Equal(t, "<nil>\n", logger.GetContents())
	})
	t.Run("error", func(t *testing.T) {
		mfs := fs.New()
		layer := new(mockLayer)
		logger := mocklog.New()
		s := &HTTPServer{
			logger: logger,
			layer:  layer,
			fs:     mfs,
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
		mfs := fs.New()
		s := &HTTPServer{
			fs: mfs,
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

func TestServeHTTP(t *testing.T) {
	s := &HTTPServer{}
	testFunc := func() {
		s.ServeHTTP(nil, nil)
	}
	assert.NotPanics(t, testFunc)
}
