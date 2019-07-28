package network

import (
	"fmt"
	"net/http"

	"golang.org/x/tools/godoc/vfs"
)

// HTTPServer is an HTTP Server
type HTTPServer struct {
	Address string
	fs      vfs.FileSystem
	net     Layer
}

// Layer is a network on which to listen and serve HTTP
type Layer interface {
	ListenAndServe(addr string, handler http.Handler) error
	ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error
}

// NewHTTPServer creates a new HTTP Server
func NewHTTPServer(fileSystem vfs.FileSystem, net Layer) *HTTPServer {
	server := &HTTPServer{
		Address: "localhost",
		fs:      fileSystem,
		net:     net,
	}
	return server
}

// Start starts the server
func (s *HTTPServer) Start() error {
	http.ListenAndServe("localhost", s)
	return fmt.Errorf("http server closed unexpectedly")
}

// ServeHTTP serves HTTP
func (s *HTTPServer) ServeHTTP(http.ResponseWriter, *http.Request) {
	s.getLayer().ListenAndServe(s.getAddress(), s)
}

func (s *HTTPServer) getFs() vfs.FileSystem {
	if s == nil {
		return nil
	}
	return s.fs
}

func (s *HTTPServer) getLayer() Layer {
	if s == nil {
		return nil
	}
	return s.net
}

func (s *HTTPServer) getAddress() string {
	if s == nil {
		return ""
	}
	return s.Address
}
