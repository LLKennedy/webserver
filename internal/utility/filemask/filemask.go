package filemask

import (
	"os"

	"golang.org/x/tools/godoc/vfs"
)

type filemask struct {
	innerSystem  vfs.FileSystem
	relativePath string
}

// Wrap encloses a file system in a new file system with a different root node
func Wrap(fs vfs.FileSystem, relativePath string) vfs.FileSystem {
	return &filemask{
		innerSystem:  fs,
		relativePath: relativePath,
	}
}

// Open opens a file on the inner file system
func (f *filemask) Open(path string) (vfs.ReadSeekCloser, error) {
	if path == "/" || path == "" {
		path = "/index.html"
	}
	fs, rpath := f.getProps(path)
	return fs.Open(rpath)
}

// Lstat gets stats on a path on the inner file system
func (f *filemask) Lstat(path string) (os.FileInfo, error) {
	fs, rpath := f.getProps(path)
	return fs.Lstat(rpath)
}

// Stat gets stats on a path on the inner file system
func (f *filemask) Stat(path string) (os.FileInfo, error) {
	fs, rpath := f.getProps(path)
	return fs.Stat(rpath)
}

// ReadDir walks a path on the inner file system
func (f *filemask) ReadDir(path string) ([]os.FileInfo, error) {
	fs, rpath := f.getProps(path)
	return fs.ReadDir(rpath)
}

// RootType gets the root type of the path on the inner file system
func (f *filemask) RootType(path string) vfs.RootType {
	fs, rpath := f.getProps(path)
	return fs.RootType(rpath)
}

// String converts the inner file system to a string
func (f *filemask) String() string {
	fs, _ := f.getProps("")
	return fs.String()
}

func (f *filemask) getProps(path string) (vfs.FileSystem, string) {
	fs := f.getFileSystem()
	rpath := f.getRelativePath(path)
	if fs == nil {
		fs = vfs.OS(".")
	}
	return fs, rpath
}

func (f *filemask) getFileSystem() vfs.FileSystem {
	if f == nil {
		return nil
	}
	return f.innerSystem
}

func (f *filemask) getRelativePath(ext string) string {
	if f == nil {
		return "" + ext
	}
	return f.relativePath + ext
}
