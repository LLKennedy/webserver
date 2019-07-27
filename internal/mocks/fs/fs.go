package fs

import (
	"os"

	"github.com/stretchr/testify/mock"
	"golang.org/x/tools/godoc/vfs"
)

type mockFS struct {
	mock.Mock
}

// New creates a new mock file system
func New() vfs.FileSystem {
	return new(mockFS)
}

// Open opens a file
func (m *mockFS) Open(path string) (vfs.ReadSeekCloser, error) {
	args := m.Called(path)
	return args.Get(0).(vfs.ReadSeekCloser), args.Error(1)
}

// Lstat gets stats
func (m *mockFS) Lstat(path string) (os.FileInfo, error) {
	args := m.Called(path)
	return args.Get(0).(os.FileInfo), args.Error(1)
}

// Stat gets stats
func (m *mockFS) Stat(path string) (os.FileInfo, error) {
	args := m.Called(path)
	return args.Get(0).(os.FileInfo), args.Error(1)
}

// ReadDir reads a directory
func (m *mockFS) ReadDir(path string) ([]os.FileInfo, error) {
	args := m.Called(path)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

// RootType gets the root type
func (m *mockFS) RootType(path string) vfs.RootType {
	return m.Called(path).Get(0).(vfs.RootType)
}

// String gets the file system as a string
func (m *mockFS) String() string {
	return m.Called().String(0)
}
