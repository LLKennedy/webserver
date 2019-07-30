package network

import (
	"net/http"
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/network"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/tools/godoc/vfs"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) ServeHTTP(http.ResponseWriter, *http.Request) {

}

func TestListenAndServe(t *testing.T) {
	handler := new(mockHandler)
	addr := "localhost"
	layer := HTTP{}
	err := layer.ListenAndServe(addr, handler)
	assert.EqualError(t, err, "listen tcp: address localhost: missing port in address")
}

func TestListenAndServeTLS(t *testing.T) {
	handler := new(mockHandler)
	addr := "localhost"
	layer := HTTP{}
	err := layer.ListenAndServeTLS(addr, "", "", handler)
	assert.EqualError(t, err, "listen tcp: address localhost: missing port in address")
}

func TestFileServer(t *testing.T) {
	mfs := vfs.NameSpace{}
	vdir := network.NewDir(mfs)
	layer := HTTP{}
	handler := layer.FileServer(vdir)
	assert.Nil(t, http.FileServer(vdir), handler)
}
