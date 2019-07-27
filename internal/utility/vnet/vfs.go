package vnet

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/tools/godoc/vfs"
)

// VDir is a virtual directory for net/http
type VDir struct {
	fs vfs.FileSystem
}

// vfile is a virtual http.File
type vfile struct {
	path string
	file vfs.ReadSeekCloser
	fs   vfs.FileSystem
}

// NewVDir returns a new virtual directory on the specified file system
func NewVDir(fs vfs.FileSystem) (*VDir, error) {
	v := &VDir{}
	_ = http.FileSystem(v)
	return v, nil
}

// Open opens a file at the specified path
func (v *VDir) Open(path string) (http.File, error) {
	if v.getFs() == nil {
		return nil, fmt.Errorf("cannot open file on nil file system")
	}
	f := &vfile{
		path: path,
		fs:   v.getFs(),
	}
	var err error
	f.file, err = v.getFs().Open(path)
	return f, err
}

// Seek seeks on the file
func (f *vfile) Seek(offset int64, whence int) (int64, error) {
	if f.getFile() == nil {
		return 0, fmt.Errorf("cannot seek on nil file")
	}
	return f.file.Seek(offset, whence)
}

// Read reads from the file
func (f *vfile) Read(p []byte) (n int, err error) {
	if f.getFile() == nil {
		return 0, fmt.Errorf("cannot read on nil file")
	}
	return f.file.Read(p)
}

// Close closes the file
func (f *vfile) Close() error {
	if f.getFile() == nil {
		return nil
	}
	return f.file.Close()
}

// Readdir reads the directory on the file system
func (f *vfile) Readdir(count int) ([]os.FileInfo, error) {
	if f.getFs() == nil {
		return nil, fmt.Errorf("cannot read directory on nil file system")
	}
	list, err := f.fs.ReadDir(f.getPath())
	if count > 0 && count < len(list) {
		list = list[:count]
	}
	return list, err
}

// Stat reads the file stats on the file
func (f *vfile) Stat() (os.FileInfo, error) {
	if f.getFs() == nil {
		return nil, fmt.Errorf("cannot get stats on nil file")
	}
	return f.fs.Stat(f.path)
}

func (v *VDir) getFs() vfs.FileSystem {
	if v == nil {
		return nil
	}
	return v.fs
}

func (f *vfile) getPath() string {
	if f == nil {
		return ""
	}
	return f.path
}

func (f *vfile) getFs() vfs.FileSystem {
	if f == nil {
		return nil
	}
	return f.fs
}

func (f *vfile) getFile() vfs.ReadSeekCloser {
	if f == nil {
		return nil
	}
	return f.file
}
