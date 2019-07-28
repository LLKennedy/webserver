package network

import (
	"fmt"
	"net/http"

	"github.com/LLKennedy/webserver/internal/mocks/vnet"
	"github.com/LLKennedy/webserver/internal/utility/logs"
	"golang.org/x/tools/godoc/vfs"
)

// HTTPServer is an HTTP Server
type HTTPServer struct {
	Address    string
	layer      Layer
	logger     logs.Logger
	fileServer http.Handler
}

// Layer is a network on which to listen and serve HTTP
type Layer interface {
	ListenAndServe(addr string, handler http.Handler) error
	ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error
}

// NewHTTPServer creates a new HTTP Server
func NewHTTPServer(logger logs.Logger, fileSystem vfs.FileSystem, layer Layer) *HTTPServer {
	server := &HTTPServer{
		logger:     logger,
		Address:    "localhost",
		layer:      layer,
		fileServer: http.FileServer(vnet.NewDir(fileSystem)),
	}
	return server
}

// Start starts the server
func (s *HTTPServer) Start() error {
	err := s.getLayer().ListenAndServe(s.getAddress(), s)
	if err != nil {
		err = fmt.Errorf("http server closed unexpectedly: %v", err)
	}
	s.getLogger().Printf("%v\n", err)
	return err
}

// ServeHTTP serves HTTP
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fileServer := s.getFs()
	fileServer.ServeHTTP(writer, request)
}

func (s *HTTPServer) getLogger() logs.Logger {
	if s == nil {
		return nil
	}
	return s.logger
}

func (s *HTTPServer) getFs() http.Handler {
	if s == nil {
		return nil
	}
	return s.fileServer
}

func (s *HTTPServer) getLayer() Layer {
	if s == nil {
		return nil
	}
	return s.layer
}

func (s *HTTPServer) getAddress() string {
	if s == nil {
		return ""
	}
	return s.Address
}
