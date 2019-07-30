package network

import (
	"net/http"
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/mocknetwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/tools/godoc/vfs"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	args := m.Called(r)
	writeFunc := args.Get(0).(func(http.ResponseWriter))
	writeFunc(w)
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
	vdir := mocknetwork.NewDir(mfs)
	layer := HTTP{}
	handler := layer.FileServer(vdir)
	assert.Nil(t, http.FileServer(vdir), handler)
}
