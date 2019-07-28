package network

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
