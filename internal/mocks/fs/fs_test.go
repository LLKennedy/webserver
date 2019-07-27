package fs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/tools/godoc/vfs"
)

func TestNewFile(t *testing.T) {
	f := NewFile("filename.ext", []byte("hello"), fmt.Errorf("an error"), fmt.Errorf("another error"), true)
	expected := new(MockFile)
	expected.name = "filename.ext"
	expected.data = []byte("hello")
	expected.err = fmt.Errorf("an error")
	assert.Equal(t, 1, len(f.ExpectedCalls))
	err := f.Close()
	assert.EqualError(t, err, "another error")
	f.Calls = []mock.Call(nil)
	f.ExpectedCalls = []*mock.Call(nil)
	assert.Equal(t, expected, f)
	f.AssertExpectations(t)
}

func TestNew(t *testing.T) {
	t.Run("no files", func(t *testing.T) {
		m := New()
		assert.Equal(t, &MockFS{}, m)
		m.AssertExpectations(t)
	})
	t.Run("one file", func(t *testing.T) {
		f := NewFile("filename.ext", []byte("hello"), fmt.Errorf("an error"), nil, false)
		m := New(f)
		exFs := new(MockFS)
		assert.Equal(t, 1, len(m.ExpectedCalls))
		f2, err := m.Open("filename.ext")
		assert.Equal(t, f, f2)
		assert.EqualError(t, err, "an error")
		m.Calls = []mock.Call(nil)
		m.ExpectedCalls = []*mock.Call(nil)
		assert.Equal(t, exFs, m)
		m.AssertExpectations(t)
	})
}

func TestOpen(t *testing.T) {
	m := New()
	m.On("Open", "filename.ext").Return(vfs.ReadSeekCloser(nil), nil)
	file, err := m.Open("filename.ext")
	assert.Nil(t, file)
	assert.NoError(t, err)
}

func TestLstat(t *testing.T) {
	m := New()
	m.On("Lstat", "filename.ext").Return(nil, nil)
	stat, err := m.Lstat("filename.ext")
	assert.Nil(t, stat)
	assert.NoError(t, err)
}

func TestStat(t *testing.T) {
	m := New()
	m.On("Stat", "filename.ext").Return(nil, nil)
	stat, err := m.Stat("filename.ext")
	assert.Nil(t, stat)
	assert.NoError(t, err)
}

func TestReadDir(t *testing.T) {
	m := New()
	m.On("ReadDir", "filename.ext").Return(nil, nil)
	files, err := m.ReadDir("filename.ext")
	assert.Nil(t, files)
	assert.NoError(t, err)
}

func TestRootType(t *testing.T) {
	m := New()
	m.On("RootType", "filename.ext").Return(vfs.RootType(""))
	rtype := m.RootType("filename.ext")
	assert.Empty(t, rtype)
}

func TestString(t *testing.T) {
	m := New()
	m.On("String").Return("")
	s := m.String()
	assert.Empty(t, s)
}
