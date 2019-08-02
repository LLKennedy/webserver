package network

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/mocklog"
	"github.com/LLKennedy/webserver/internal/mocks/mocknetwork"
	"github.com/LLKennedy/webserver/internal/utility/filemask"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/godoc/vfs"
)

func TestNewHTTPServer(t *testing.T) {
	mfs := vfs.NewNameSpace()
	layer := HTTP{}
	logger := mocklog.New()
	s := NewHTTPServer(logger, mfs, layer)
	assert.Equal(t, &HTTPServer{logger: logger, Address: "localhost", Port: "80", layer: layer, fileServer: http.FileServer(mocknetwork.NewDir(filemask.Wrap(mfs, "build"))), staticServer: http.FileServer(mocknetwork.NewDir(filemask.Wrap(mfs, "build/static")))}, s)
	assert.Equal(t, "", logger.GetContents())
}

func TestStart(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		mfs := vfs.NewNameSpace()
		layer := new(mocknetwork.Layer)
		logger := mocklog.New()
		s := &HTTPServer{
			logger:       logger,
			layer:        layer,
			fileServer:   http.FileServer(mocknetwork.NewDir(filemask.Wrap(mfs, "build"))),
			staticServer: http.FileServer(mocknetwork.NewDir(filemask.Wrap(mfs, "build/static"))),
		}
		layer.On("ListenAndServe", fmt.Sprintf("%s:%s", s.getAddress(), s.getPort()), s).Return(nil)
		err := s.Start()
		assert.NoError(t, err)
		assert.Equal(t, "<nil>\n", logger.GetContents())
	})
	t.Run("error", func(t *testing.T) {
		mfs := vfs.NewNameSpace()
		layer := new(mocknetwork.Layer)
		logger := mocklog.New()
		s := &HTTPServer{
			logger:       logger,
			layer:        layer,
			fileServer:   http.FileServer(mocknetwork.NewDir(filemask.Wrap(mfs, "build"))),
			staticServer: http.FileServer(mocknetwork.NewDir(filemask.Wrap(mfs, "build/static"))),
		}
		layer.On("ListenAndServe", fmt.Sprintf("%s:%s", s.getAddress(), s.getPort()), s).Return(fmt.Errorf("some network error"))
		err := s.Start()
		assert.EqualError(t, err, "http server closed unexpectedly: some network error")
		assert.Equal(t, "http server closed unexpectedly: some network error\n", logger.GetContents())
	})
}

func TestGetFs(t *testing.T) {
	t.Run("nil server", func(t *testing.T) {
		var s *HTTPServer
		ffs, sfs := s.getFs()
		assert.Nil(t, ffs)
		assert.Nil(t, sfs)
	})
	t.Run("non-nil server", func(t *testing.T) {
		mfs := new(mockHandler)
		s := &HTTPServer{
			fileServer:   mfs,
			staticServer: mfs,
		}
		ffs, sfs := s.getFs()
		assert.Equal(t, mfs, ffs)
		assert.Equal(t, mfs, sfs)
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
		mlayer := new(mocknetwork.Layer)
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

type mockLogger struct {
	fatalCount int
	buf        bytes.Buffer
}

func (m *mockLogger) Println(v ...interface{}) {
	m.buf.WriteString(fmt.Sprintln(v...))
}

func (m *mockLogger) Printf(format string, v ...interface{}) {
	m.buf.WriteString(fmt.Sprintf(format, v...))
}

func (m *mockLogger) Fatalf(format string, v ...interface{}) {
	m.fatalCount++
	m.buf.WriteString(fmt.Sprintf(format, v...))
}

func TestServeHTTP(t *testing.T) {
	mfs := vfs.NewNameSpace()
	s := &HTTPServer{
		fileServer: http.FileServer(mocknetwork.NewDir(mfs)),
		logger:     new(mockLogger),
	}
	testFunc := func() {
		defer func() {
			if r := recover(); r != nil {
				assert.Failf(t, "caught panic", "%v\n%s", r, debug.Stack())
			}
		}()
		s.ServeHTTP(new(mockResponseWriter), &http.Request{URL: &url.URL{}, Host: "", RemoteAddr: ""})
	}
	assert.NotPanics(t, testFunc)
}
