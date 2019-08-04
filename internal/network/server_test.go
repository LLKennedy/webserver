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
	"time"

	"github.com/LLKennedy/webserver/internal/mocks/fs"
	"github.com/LLKennedy/webserver/internal/mocks/mocklog"
	"github.com/LLKennedy/webserver/internal/mocks/mocknetwork"
	"github.com/LLKennedy/webserver/internal/utility/config"
	"github.com/LLKennedy/webserver/internal/utility/filemask"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/tools/godoc/vfs"
)

func TestNewHTTPServer(t *testing.T) {
	mfs := fs.New()
	mfs.On("Open", mock.Anything).Return(fs.NewFile("", []byte("{}"), nil, nil, true), nil)
	layer := HTTP{}
	logger := mocklog.New()
	s := NewHTTPServer(logger, mfs, layer)
	assert.Equal(t, &HTTPServer{
		Options:    config.DefaultOptions(),
		logger:     logger,
		layer:      layer,
		fileSystem: mfs,
		fileServer: http.FileServer(mocknetwork.NewDir(filemask.Wrap(mfs, "build"))),
	}, s)
	assert.Equal(t, "", logger.GetContents())
}

func catchPanic(t *testing.T) {
	if r := recover(); r != nil {
		assert.Failf(t, "caught panic", "%v\n%s", r, debug.Stack())
	}
}

func TestStart(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		defer catchPanic(t)
		testHash := "'sha256-5As4+3YpY62+l38PsxCEkjB1R4YtyktBtRScTJ3fyLU='"
		mfs := fs.New(fs.NewFile("build/index.html", []byte(testHash), nil, nil, true))
		layer := new(mocknetwork.Layer)
		logger := mocklog.New()
		s := &HTTPServer{
			Options:    config.DefaultOptions(),
			logger:     logger,
			layer:      layer,
			fileSystem: mfs,
			fileServer: http.FileServer(mocknetwork.NewDir(filemask.Wrap(mfs, "build"))),
		}
		expectedInsecure := &insecureServer{Address: config.DefaultOptions().Address, scriptHash: testHash}
		layer.On("ListenAndServe", fmt.Sprintf("%s:%d", s.getOptions().Address, s.getOptions().InsecurePort), expectedInsecure).Return(nil)
		err := s.Start()
		assert.NoError(t, err)
		assert.Equal(t, "", logger.GetContents())
	})
	t.Run("error getting script hash", func(t *testing.T) {
		defer catchPanic(t)
		mfs := fs.New(fs.NewFile("build/index.html", nil, fmt.Errorf("can't open file"), nil, true))
		layer := new(mocknetwork.Layer)
		logger := mocklog.New()
		s := &HTTPServer{
			logger:     logger,
			layer:      layer,
			fileSystem: mfs,
			fileServer: http.FileServer(mocknetwork.NewDir(filemask.Wrap(mfs, "build"))),
		}
		layer.On("ListenAndServe", fmt.Sprintf("%s:%d", s.getOptions().Address, s.getOptions().Port), s).Return(fmt.Errorf("some network error"))
		err := s.Start()
		assert.EqualError(t, err, "could not read script hash: could not open index file: can't open file")
		assert.Equal(t, "could not open index file: can't open file\ncould not read script hash: could not open index file: can't open file\n", logger.GetContents())
	})
	t.Run("error starting HTTP server", func(t *testing.T) {
		defer catchPanic(t)
		mfs := fs.New(fs.NewFile("build/index.html", []byte("'sha256-5As4+3YpY62+l38PsxCEkjB1R4YtyktBtRScTJ3fyLU='"), nil, nil, true))
		layer := new(mocknetwork.Layer)
		logger := mocklog.New()
		s := &HTTPServer{
			logger:     logger,
			layer:      layer,
			fileSystem: mfs,
			fileServer: http.FileServer(mocknetwork.NewDir(filemask.Wrap(mfs, "build"))),
		}
		layer.On("ListenAndServe", fmt.Sprintf("%s:%d", s.getOptions().Address, s.getOptions().Port), s).Return(fmt.Errorf("some network error"))
		err := s.Start()
		assert.EqualError(t, err, "http server closed unexpectedly: some network error")
		assert.Equal(t, "http server closed unexpectedly: some network error\n", logger.GetContents())
	})
}

func TestReadScriptHash(t *testing.T) {
	t.Run("fail on file open", func(t *testing.T) {
		mfs := fs.New(fs.NewFile("build/index.html", nil, fmt.Errorf("some error"), nil, false))
		logger := new(mockLogger)
		s := &HTTPServer{
			fileSystem: mfs,
			logger:     logger,
		}
		hash, err := s.readScriptHash()
		assert.Empty(t, hash)
		assert.EqualError(t, err, "could not open index file: some error")
	})
	t.Run("fail on regex", func(t *testing.T) {
		mfs := fs.New(fs.NewFile("build/index.html", []byte{}, nil, nil, true))
		logger := new(mockLogger)
		s := &HTTPServer{
			fileSystem: mfs,
			logger:     logger,
		}
		hash, err := s.readScriptHash()
		assert.Empty(t, hash)
		assert.EqualError(t, err, "could not find script hash in index file")
	})
	t.Run("fail on navigation", func(t *testing.T) {
		indexFile := fs.NewFile("build/index.html", []byte("'sha256-5As4+3YpY62+l38PsxCEkjB1R4YtyktBtRScTJ3fyLU='"), nil, nil, true)
		indexFile.OverrideBuffer = true
		indexFile.On("Seek", int64(0), 0).Return(int64(0), fmt.Errorf("seek error"))
		mfs := fs.New(indexFile)
		logger := new(mockLogger)
		s := &HTTPServer{
			fileSystem: mfs,
			logger:     logger,
		}
		hash, err := s.readScriptHash()
		assert.Empty(t, hash)
		assert.EqualError(t, err, "could not navigate to specified location in index file")
	})
	t.Run("successful extraction", func(t *testing.T) {
		indexFile := fs.NewFile("build/index.html", []byte("'sha256-5As4+3YpY62+l38PsxCEkjB1R4YtyktBtRScTJ3fyLU='"), nil, nil, true)
		mfs := fs.New(indexFile)
		logger := new(mockLogger)
		s := &HTTPServer{
			fileSystem: mfs,
			logger:     logger,
		}
		hash, err := s.readScriptHash()
		assert.Equal(t, "'sha256-5As4+3YpY62+l38PsxCEkjB1R4YtyktBtRScTJ3fyLU='", hash)
		assert.NoError(t, err)
	})
}

func TestGetFs(t *testing.T) {
	t.Run("nil server", func(t *testing.T) {
		defer catchPanic(t)
		var s *HTTPServer
		ffs := s.getFs()
		assert.Nil(t, ffs)
	})
	t.Run("non-nil server", func(t *testing.T) {
		defer catchPanic(t)
		mfs := new(mockHandler)
		s := &HTTPServer{
			fileServer: mfs,
		}
		ffs := s.getFs()
		assert.Equal(t, mfs, ffs)
	})
}

func TestGetLogger(t *testing.T) {
	t.Run("nil server", func(t *testing.T) {
		defer catchPanic(t)
		var s *HTTPServer
		logger := s.getLogger()
		assert.Nil(t, logger)
	})
	t.Run("non-nil server", func(t *testing.T) {
		defer catchPanic(t)
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
		defer catchPanic(t)
		var s *HTTPServer
		layer := s.getLayer()
		assert.Nil(t, layer)
	})
	t.Run("non-nil server", func(t *testing.T) {
		defer catchPanic(t)
		mlayer := new(mocknetwork.Layer)
		s := &HTTPServer{
			layer: mlayer,
		}
		layer := s.getLayer()
		assert.Equal(t, mlayer, layer)
	})
}

func TestGetScriptHash(t *testing.T) {
	t.Run("nil server", func(t *testing.T) {
		defer catchPanic(t)
		var s *HTTPServer
		layer := s.getScriptHash()
		assert.Empty(t, layer)
	})
	t.Run("non-nil server", func(t *testing.T) {
		defer catchPanic(t)
		inaddr := "12"
		s := &HTTPServer{
			scriptHash: inaddr,
		}
		outaddr := s.getScriptHash()
		assert.Equal(t, inaddr, outaddr)
	})
}

func TestGetFileSystem(t *testing.T) {
	t.Run("nil server", func(t *testing.T) {
		defer catchPanic(t)
		defaultFs := vfs.OS(".")
		var s *HTTPServer
		rfs := s.getFileSystem()
		assert.Equal(t, defaultFs, rfs)
	})
	t.Run("non-nil server", func(t *testing.T) {
		defer catchPanic(t)
		mfs := fs.New()
		s := &HTTPServer{
			fileSystem: mfs,
		}
		rfs := s.getFileSystem()
		assert.Equal(t, mfs, rfs)
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
	t.Run("root file", func(t *testing.T) {
		defer catchPanic(t)
		rootFile := fs.NewFile("/", []byte(""), nil, nil, true)
		rootFile.On("IsDir").Return(false)
		mfs := fs.New(rootFile)
		mfs.On("Stat", "/").Return(rootFile, nil)
		s := &HTTPServer{
			fileServer: http.FileServer(mocknetwork.NewDir(mfs)),
			logger:     new(mockLogger),
		}
		s.ServeHTTP(new(mockResponseWriter), &http.Request{URL: &url.URL{Path: "/"}, Host: "", RemoteAddr: ""})
		mfs.AssertExpectations(t)
		rootFile.AssertExpectations(t)
	})
	t.Run("js file", func(t *testing.T) {
		defer catchPanic(t)
		scriptFile := fs.NewFile("/something.js", []byte(""), nil, nil, true)
		scriptFile.On("IsDir").Return(false)
		scriptFile.On("Name").Return("something.js")
		scriptFile.On("ModTime").Return(time.Now())
		scriptFile.On("Size").Return(int64(len("")))
		scriptFile.On("IsDir").Return(false)
		mfs := fs.New(scriptFile)
		mfs.On("Stat", "/something.js").Return(scriptFile, nil)
		s := &HTTPServer{
			fileServer: http.FileServer(mocknetwork.NewDir(mfs)),
			logger:     new(mockLogger),
		}
		s.ServeHTTP(new(mockResponseWriter), &http.Request{RequestURI: "something.js", URL: &url.URL{Path: "/something.js"}, Host: "", RemoteAddr: ""})
		mfs.AssertExpectations(t)
		scriptFile.AssertExpectations(t)
	})
	t.Run("ts file", func(t *testing.T) {
		defer catchPanic(t)
		scriptFile := fs.NewFile("/something.ts", []byte(""), nil, nil, true)
		scriptFile.On("IsDir").Return(false)
		scriptFile.On("Name").Return("something.ts")
		scriptFile.On("ModTime").Return(time.Now())
		scriptFile.On("Size").Return(int64(len("")))
		scriptFile.On("IsDir").Return(false)
		mfs := fs.New(scriptFile)
		mfs.On("Stat", "/something.ts").Return(scriptFile, nil)
		s := &HTTPServer{
			fileServer: http.FileServer(mocknetwork.NewDir(mfs)),
			logger:     new(mockLogger),
		}
		s.ServeHTTP(new(mockResponseWriter), &http.Request{RequestURI: "something.ts", URL: &url.URL{Path: "/something.ts"}, Host: "", RemoteAddr: ""})
		mfs.AssertExpectations(t)
		scriptFile.AssertExpectations(t)
	})
}
