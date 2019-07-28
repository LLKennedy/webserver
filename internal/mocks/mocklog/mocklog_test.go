package mocklog

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	m := New()
	assert.Equal(t, &MockLogger{buffer: bytes.NewBuffer(nil)}, m)
}

func TestPrintln(t *testing.T) {
	m := New()
	m.Println("hello")
	eol := fmt.Sprintln("")
	combined := "hello" + eol
	assert.Equal(t, bytes.NewBuffer([]byte(combined)), m.GetBuffer())
}

func TestPrintf(t *testing.T) {
	m := New()
	m.Printf("hello")
	assert.Equal(t, bytes.NewBuffer([]byte("hello")), m.GetBuffer())
}

func TestFatalf(t *testing.T) {
	m := New()
	m.Fatalf("hello")
	assert.Equal(t, bytes.NewBuffer([]byte("hello")), m.GetBuffer())
	assert.Equal(t, 1, m.GetFatalCount())
}

func TestGetBuffer(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var m *MockLogger
		buf := m.GetBuffer()
		assert.Nil(t, buf)
	})
	t.Run("non-nil", func(t *testing.T) {
		m := New()
		buf := m.GetBuffer()
		assert.Equal(t, bytes.NewBuffer(nil), buf)
	})
}

func TestGetFatalCount(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var m *MockLogger
		count := m.GetFatalCount()
		assert.Equal(t, 0, count)
	})
	t.Run("non-nil", func(t *testing.T) {
		m := New()
		count := m.GetFatalCount()
		assert.Equal(t, 0, count)
	})
}

func TestSetFatalCount(t *testing.T) {
	m := New()
	m.SetFatalCount(2)
	assert.Equal(t, 2, m.GetFatalCount())
}

func TestGetContents(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var m *MockLogger
		contents := m.GetContents()
		assert.Equal(t, "", contents)
	})
	t.Run("non-nil", func(t *testing.T) {
		m := New()
		contents := m.GetContents()
		assert.Equal(t, "", contents)
	})
}
