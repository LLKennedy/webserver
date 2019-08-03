package fs

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/tools/godoc/vfs"
)

func TestNewFile(t *testing.T) {
	f := NewFile("filename.ext", []byte("hello"), fmt.Errorf("an error"), fmt.Errorf("another error"), true)
	expected := new(MockFile)
	expected.name = "filename.ext"
	expected.data = []byte("hello")
	expected.Buf = bytes.NewReader([]byte("hello"))
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
	t.Run("no error", func(t *testing.T) {
		f := NewFile("filename.ext", []byte("hello"), nil, nil, false)
		sought, err := f.Seek(2, 1)
		assert.Equal(t, sought, int64(2))
		assert.NoError(t, err)
	})
	t.Run("error", func(t *testing.T) {
		f := NewFile("filename.ext", []byte("hello"), nil, nil, false)
		sought, err := f.Seek(50, 100)
		assert.Equal(t, sought, int64(0))
		assert.EqualError(t, err, "bytes.Reader.Seek: invalid whence")
	})
	t.Run("override", func(t *testing.T) {
		f := NewFile("filename.ext", []byte("hello"), nil, nil, false)
		f.OverrideBuffer = true
		f.On("Seek", int64(50), 100).Return(int64(12), fmt.Errorf("some error"))
		sought, err := f.Seek(50, 100)
		assert.Equal(t, int64(12), sought)
		assert.EqualError(t, err, "some error")
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

func TestName(t *testing.T) {
	f := new(MockFile)
	f.On("Name").Return("somefile.ext")
	assert.Equal(t, "somefile.ext", f.Name())
	f.AssertExpectations(t)
}

func TestSize(t *testing.T) {
	f := new(MockFile)
	f.On("Size").Return(int64(45))
	assert.Equal(t, int64(45), f.Size())
	f.AssertExpectations(t)
}

func TestMode(t *testing.T) {
	f := new(MockFile)
	f.On("Mode").Return(os.FileMode(1))
	assert.Equal(t, os.FileMode(1), f.Mode())
	f.AssertExpectations(t)
}

func TestModTime(t *testing.T) {
	now := time.Now()
	f := new(MockFile)
	f.On("ModTime").Return(now)
	assert.Equal(t, now, f.ModTime())
	f.AssertExpectations(t)
}

func TestIsDir(t *testing.T) {
	f := new(MockFile)
	f.On("IsDir").Return(true)
	assert.True(t, f.IsDir())
	f.AssertExpectations(t)
}

func TestSys(t *testing.T) {
	f := new(MockFile)
	f.On("Sys").Return(nil)
	assert.Nil(t, f.Sys())
	f.AssertExpectations(t)
}
