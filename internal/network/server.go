package network

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/LLKennedy/webserver/internal/mocks/mocknetwork"
	"github.com/LLKennedy/webserver/internal/utility/config"
	"github.com/LLKennedy/webserver/internal/utility/filemask"
	"github.com/LLKennedy/webserver/internal/utility/logs"
	"golang.org/x/tools/godoc/vfs"
)

var jsFile = regexp.MustCompile(`\.js$`)
var tsFile = regexp.MustCompile(`\.ts$`)

// HTTPServer is an HTTP Server
type HTTPServer struct {
	Options        config.Options
	scriptHash     string
	secure         Layer
	insecure       Layer
	logger         logs.Logger
	fileSystem     vfs.FileSystem
	fileServer     http.Handler
	insecureServer *insecureServer
}

type insecureServer struct {
	Address    string
	scriptHash string
}

// Layer is a network on which to listen and serve HTTP
type Layer interface {
	Shutdown(ctx context.Context) error
	ListenAndServe(addr string, handler http.Handler) error
	ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error
	FileServer(http.FileSystem) http.Handler
}

// NewHTTPServer creates a new HTTP Server
func NewHTTPServer(logger logs.Logger, fileSystem vfs.FileSystem, secure, insecure Layer) *HTTPServer {
	options, err := config.Load(fileSystem, os.Args)
	server := &HTTPServer{
		Options:    options,
		logger:     logger,
		secure:     secure,
		insecure:   insecure,
		fileSystem: fileSystem,
		fileServer: http.FileServer(mocknetwork.NewDir(filemask.Wrap(fileSystem, "build"))),
	}
	if err != nil {
		server.getLogger().Printf("problem getting config: %v", err)
	}
	return server
}

type runeVFS struct {
	file vfs.ReadSeekCloser
}

func (rv *runeVFS) ReadRune() (r rune, size int, err error) {
	next := make([]byte, 1)
	size, err = rv.file.Read(next)
	if err == nil {
		r = rune(next[0])
	}
	return
}

// Start starts the server
func (s *HTTPServer) Start() (err error) {
	s.scriptHash, err = s.readScriptHash()
	if err != nil {
		err = fmt.Errorf("could not read script hash: %v", err)
		s.getLogger().Printf("%v\n", err)
		return err
	}
	s.insecureServer = &insecureServer{Address: s.getOptions().Address, scriptHash: s.getScriptHash()}
	go s.getInsecure().ListenAndServe(fmt.Sprintf("%s:%d", s.getOptions().Address, s.getOptions().InsecurePort), s.insecureServer)
	err = s.getSecure().ListenAndServeTLS(fmt.Sprintf("%s:%d", s.getOptions().Address, s.getOptions().Port), s.getOptions().CertFile, s.getOptions().KeyFile, s)
	if err != nil {
		err = fmt.Errorf("http server closed unexpectedly: %v", err)
		s.getLogger().Printf("%v\n", err)
	}
	return err
}

// ServeHTTP serves HTTP
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	protocol := "https"
	setHeaders(writer, s.getOptions().Address, protocol, s.getScriptHash())
	uri := request.RequestURI
	if jsFile.MatchString(uri) || tsFile.MatchString(uri) {
		writer.Header().Set("Content-Type", "application/javascript")
	}
	s.getFs().ServeHTTP(writer, request)
}

// ServeHTTP serves HTTP
func (i *insecureServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	protocol := "http"
	addr := config.DefaultOptions().Address
	scriptHash := "'none'"
	if i != nil {
		addr = i.Address
		scriptHash = i.scriptHash
	}
	setHeaders(writer, addr, protocol, scriptHash)
	handlerFunc := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "https://"+req.Host+req.URL.Path, http.StatusMovedPermanently)
	})
	handlerFunc.ServeHTTP(writer, request)
}

func (s *HTTPServer) readScriptHash() (scriptHash string, err error) {
	hash := regexp.MustCompile(`'sha256-[a-zA-Z0-9+/=]{44}'`)
	indexFile, err := s.getFileSystem().Open("build/index.html")
	if err != nil {
		err = fmt.Errorf("could not open index file: %v", err)
		s.getLogger().Printf("%v\n", err)
		return
	}
	defer indexFile.Close()
	locs := hash.FindReaderIndex(&runeVFS{file: indexFile})
	if locs == nil || len(locs) != 2 {
		err = fmt.Errorf("could not find script hash in index file")
		s.getLogger().Printf("%v\n", err)
		return
	}
	_, err = indexFile.Seek(int64(locs[0]), io.SeekStart)
	if err != nil {
		err = fmt.Errorf("could not navigate to specified location in index file")
		s.getLogger().Printf("%v\n", err)
		return
	}
	hashBytes := make([]byte, locs[1]-locs[0])
	indexFile.Read(hashBytes) // We've already read exactly this section of the file, don't bother error-handling the same thing again
	scriptHash = string(hashBytes)
	return
}

func setHeaders(writer http.ResponseWriter, addr, protocol, scriptHash string) {
	writer.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=31536000; includeSubDomains"))
	writer.Header().Add("Content-Security-Policy", fmt.Sprintf("default-src 'self' %s", scriptHash))
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

func (s *HTTPServer) getFs() http.Handler {
	if s == nil {
		return nil
	}
	return s.fileServer
}

func (s *HTTPServer) getSecure() Layer {
	if s == nil {
		return &HTTP{}
	}
	return s.secure
}

func (s *HTTPServer) getInsecure() Layer {
	if s == nil {
		return &HTTP{}
	}
	return s.insecure
}

func (s *HTTPServer) getFileSystem() vfs.FileSystem {
	if s == nil {
		return vfs.OS(".")
	}
	return s.fileSystem
}

func (s *HTTPServer) getScriptHash() string {
	if s == nil {
		return ""
	}
	return s.scriptHash
}

func (s *HTTPServer) getOptions() config.Options {
	if s == nil {
		return config.DefaultOptions()
	}
	return s.Options
}
