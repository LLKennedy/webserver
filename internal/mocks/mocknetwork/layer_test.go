package mocknetwork

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	args := m.Called(r)
	writeFunc := args.Get(0).(func(http.ResponseWriter))
	writeFunc(w)
}

func TestShutdown(t *testing.T) {
	ctx := context.Background()
	m := new(Layer)
	m.On("Shutdown", ctx).Return(nil)
	err := m.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestListenAndServe(t *testing.T) {
	m := new(Layer)
	m.On("ListenAndServe", "someaddress:someport", nil).Return(fmt.Errorf("some error"))
	err := m.ListenAndServe("someaddress:someport", nil)
	assert.EqualError(t, err, "some error")
}

func TestListenAndServeTLS(t *testing.T) {
	m := new(Layer)
	m.On("ListenAndServeTLS", "someaddress:someport", "my.crt", "my.key", nil).Return(fmt.Errorf("some error"))
	err := m.ListenAndServeTLS("someaddress:someport", "my.crt", "my.key", nil)
	assert.EqualError(t, err, "some error")
}

func TestFileServer(t *testing.T) {
	t.Run("nil return", func(t *testing.T) {
		m := new(Layer)
		m.On("FileServer", nil).Return(nil)
		handler := m.FileServer(nil)
		assert.Nil(t, handler)
	})
	t.Run("non-nil return", func(t *testing.T) {
		m := new(Layer)
		m.On("FileServer", nil).Return(new(mockHandler))
		handler := m.FileServer(nil)
		assert.Equal(t, new(mockHandler), handler)
	})
}
