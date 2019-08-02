package network

import (
	"fmt"
	"net/http"

	"github.com/LLKennedy/webserver/internal/mocks/mocknetwork"
	"github.com/LLKennedy/webserver/internal/utility/filemask"
	"github.com/LLKennedy/webserver/internal/utility/logs"
	"golang.org/x/tools/godoc/vfs"
)

// HTTPServer is an HTTP Server
type HTTPServer struct {
	Address      string
	Port         string
	layer        Layer
	logger       logs.Logger
	fileServer   http.Handler
	staticServer http.Handler
}

// Layer is a network on which to listen and serve HTTP
type Layer interface {
	ListenAndServe(addr string, handler http.Handler) error
	ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error
	FileServer(http.FileSystem) http.Handler
}

// NewHTTPServer creates a new HTTP Server
func NewHTTPServer(logger logs.Logger, fileSystem vfs.FileSystem, layer Layer) *HTTPServer {
	server := &HTTPServer{
		logger:       logger,
		Address:      "localhost",
		Port:         "80",
		layer:        layer,
		fileServer:   http.FileServer(mocknetwork.NewDir(filemask.Wrap(fileSystem, "build"))),
		staticServer: http.FileServer(mocknetwork.NewDir(filemask.Wrap(fileSystem, "build/static"))),
	}
	return server
}

// Start starts the server
func (s *HTTPServer) Start() error {
	err := s.getLayer().ListenAndServe(fmt.Sprintf("%s:%s", s.getAddress(), s.getPort()), s)
	if err != nil {
		err = fmt.Errorf("http server closed unexpectedly: %v", err)
	}
	s.getLogger().Printf("%v\n", err)
	return err
}

// ServeHTTP serves HTTP
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	addr := s.getAddress()
	protocol := "http"
	s.getLogger().Printf("request: host=%s, remoteAddr=%s", request.Host, request.RemoteAddr)
	setHeaders(writer, addr, protocol)
	// writer.Header().Set("Feature-Policy", "*")
	head, remainder := getPathNode(request.URL.Path)
	fileServer, staticServer := s.getFs()
	if head == "static" {
		s.getLogger().Printf("static path: %s", remainder)
		staticServer.ServeHTTP(writer, request)
	} else {
		s.getLogger().Printf("other path: %s", remainder)
		fileServer.ServeHTTP(writer, request)
	}
}

func setHeaders(writer http.ResponseWriter, addr, protocol string) {
	writer.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=31536000; includeSubDomains"))
	writer.Header().Set("Content-Security-Policy", fmt.Sprintf("default-src 'self'; style-src 'self' 'nonce-yhbSSLk5nTP4sgETaQx5Lg=='; script-src 'self' 'sha256-5As4+3YpY62+l38PsxCEkjB1R4YtyktBtRScTJ3fyLU='"))
	writer.Header().Set("X-Frame-Options", fmt.Sprintf("SAMEORIGIN"))
	writer.Header().Set("X-Content-Type-Options", fmt.Sprintf("nosniff"))
	writer.Header().Set("X-XSS-Protection", fmt.Sprintf("1; mode=block; report=%s://%s/api/security/report", protocol, addr))
	writer.Header().Set("Referrer-Policy", fmt.Sprintf("no-referrer"))
	writer.Header().Set("Set-Cookie", fmt.Sprintf("HttpOnly;Secure;SameSite=Strict"))
}

func (s *HTTPServer) getLogger() logs.Logger {
	if s == nil {
		return nil
	}
	return s.logger
}

func (s *HTTPServer) getFs() (files http.Handler, static http.Handler) {
	if s == nil {
		return nil, nil
	}
	return s.fileServer, s.staticServer
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

func (s *HTTPServer) getPort() string {
	if s == nil {
		return ""
	}
	return s.Port
}
