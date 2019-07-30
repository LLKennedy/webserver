package network

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListenAndServe(t *testing.T) {
	m := new(MockLayer)
	m.On("ListenAndServe", "someaddress:someport", nil).Return(fmt.Errorf("some error"))
	err := m.ListenAndServe("someaddress:someport", nil)
	assert.EqualError(t, err, "some error")
}

func TestListenAndServeTLS(t *testing.T) {
	m := new(MockLayer)
	m.On("ListenAndServeTLS", "someaddress:someport", "my.crt", "my.key", nil).Return(fmt.Errorf("some error"))
	err := m.ListenAndServeTLS("someaddress:someport", "my.crt", "my.key", nil)
	assert.EqualError(t, err, "some error")
}

func TestFileServer(t *testing.T) {
	m := new(MockLayer)
	m.On("FileServer", nil).Return(nil)
	handler := m.FileServer(nil)
	assert.Nil(t, handler)
}
