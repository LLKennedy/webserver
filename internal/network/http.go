package network

import "net/http"

// HTTP is the net/http implementation of the Network interface
type HTTP struct{}

// ListenAndServe calls ListenAndServe on net/http
func (h HTTP) ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

// ListenAndServeTLS calls ListenAndServeTLS on net/http
func (h HTTP) ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
}

// FileServer servies the file system over HTTP
func (h HTTP) FileServer(fs http.FileSystem) http.Handler {
	return http.FileServer(fs)
}
