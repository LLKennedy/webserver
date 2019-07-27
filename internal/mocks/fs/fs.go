package fs

import (
	"os"

	"github.com/stretchr/testify/mock"
	"golang.org/x/tools/godoc/vfs"
)

// MockFS is a mock file system
type MockFS struct {
	mock.Mock
}

// New creates a new mock file system
func New() *MockFS {
	m := new(MockFS)
	_ = vfs.FileSystem(m)
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
