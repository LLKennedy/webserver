package fs

import (
	"bytes"
	"os"

	"github.com/stretchr/testify/mock"
	"golang.org/x/tools/godoc/vfs"
)

// MockFS is a mock file system
type MockFS struct {
	mock.Mock
}

// MockFile is a mock file
type MockFile struct {
	mock.Mock
	name string
	data []byte
	err  error
	buf  *bytes.Reader
}

// New creates a new mock file system
func New(files ...*MockFile) *MockFS {
	m := new(MockFS)
	for _, file := range files {
		file.buf = bytes.NewReader(file.data)
		m.On("Open", file.name).Return(file, file.err)
	}
	_ = vfs.FileSystem(m)
	return m
}

// NewFile creates a new file system entry
func NewFile(name string, data []byte, openErr, closeErr error, expectClose bool) *MockFile {
	m := new(MockFile)
	m.name = name
	m.err = openErr
	m.data = data
	if expectClose {
		m.On("Close").Return(closeErr)
	}
	_ = vfs.ReadSeekCloser(m)
	return m
}

// Open opens a file
func (m *MockFS) Open(path string) (vfs.ReadSeekCloser, error) {
	args := m.Called(path)
	var file vfs.ReadSeekCloser
	file, _ = args.Get(0).(vfs.ReadSeekCloser)
	return file, args.Error(1)
}

// Lstat gets stats
func (m *MockFS) Lstat(path string) (os.FileInfo, error) {
	args := m.Called(path)
	var file os.FileInfo
	file, _ = args.Get(0).(os.FileInfo)
	return file, args.Error(1)
}

// Stat gets stats
func (m *MockFS) Stat(path string) (os.FileInfo, error) {
	args := m.Called(path)
	var file os.FileInfo
	file, _ = args.Get(0).(os.FileInfo)
	return file, args.Error(1)
}

// ReadDir reads a directory
func (m *MockFS) ReadDir(path string) ([]os.FileInfo, error) {
	args := m.Called(path)
	var file []os.FileInfo
	file, _ = args.Get(0).([]os.FileInfo)
	return file, args.Error(1)
}

// RootType gets the root type
func (m *MockFS) RootType(path string) vfs.RootType {
	return m.Called(path).Get(0).(vfs.RootType)
}

// String gets the file system as a string
func (m *MockFS) String() string {
	return m.Called().String(0)
}

// Seek seeks on the file
func (f *MockFile) Seek(offset int64, whence int) (int64, error) {
	return f.getBuf().Seek(offset, whence)
}

// Read reads from the file
func (f *MockFile) Read(p []byte) (n int, err error) {
	return f.getBuf().Read(p)
}

// Close closes the file
func (f *MockFile) Close() error {
	args := f.Called()
	return args.Error(0)
}

func (f *MockFile) getBuf() *bytes.Reader {
	if f == nil {
		var buf *bytes.Reader
		return buf
	}
	return f.buf
}
