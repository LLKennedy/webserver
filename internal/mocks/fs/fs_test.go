package fs

import (
	"bytes"
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
	expected.buf = bytes.NewReader([]byte("hello"))
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

func TestSeek(t *testing.T) {
	f := NewFile("filename.ext", []byte("hello"), nil, nil, false)
	t.Run("no error", func(t *testing.T) {
		sought, err := f.Seek(2, 1)
		assert.Equal(t, sought, int64(2))
		assert.NoError(t, err)
	})
	t.Run("error", func(t *testing.T) {
		sought, err := f.Seek(50, 100)
		assert.Equal(t, sought, int64(0))
		assert.EqualError(t, err, "bytes.Reader.Seek: invalid whence")
	})
}

func TestRead(t *testing.T) {
	f := NewFile("filename.ext", []byte("hello"), nil, nil, false)
	buf := make([]byte, 20)
	bytesRead, err := f.Read(buf)
	assert.Equal(t, 5, bytesRead)
	assert.NoError(t, err)
	assert.Equal(t, append([]byte("hello"), make([]byte, 15)...), buf)
}

func TestClose(t *testing.T) {
	f := NewFile("filename.ext", []byte("hello"), nil, fmt.Errorf("error closing"), true)
	err := f.Close()
	assert.EqualError(t, err, "error closing")
}

func TestGetBuf(t *testing.T) {
	var f *MockFile
	buf := f.getBuf()
	assert.Nil(t, buf)
}
