package mocknetwork

import (
	"fmt"
	"os"
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/fs"
	"github.com/stretchr/testify/assert"
)

func TestNewDir(t *testing.T) {
	mfs := fs.New()
	v := NewDir(mfs)
	expected := new(VDir)
	expected.fs = mfs
	assert.Equal(t, expected, v)
}

func TestOpen(t *testing.T) {
	t.Run("non-nil directory", func(t *testing.T) {
		f := fs.NewFile("filepath.ext", []byte("hello"), fmt.Errorf("somehow opened"), nil, false)
		mfs := fs.New(f)
		v := NewDir(mfs)
		file, err := v.Open("filepath.ext")
		assert.EqualError(t, err, "somehow opened")
		vf, ok := file.(*vfile)
		if assert.True(t, ok) {
			assert.Equal(t, f, vf.file)
			assert.Equal(t, mfs, vf.fs)
			assert.Equal(t, "filepath.ext", vf.path)
		}
	})
	t.Run("nil directory", func(t *testing.T) {
		var v *VDir
		file, err := v.Open("anything")
		assert.Nil(t, file)
		assert.EqualError(t, err, "cannot open file on nil file system")
	})
}

func TestSeek(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		var file *vfile
		sought, err := file.Seek(10, 10)
		assert.Equal(t, int64(0), sought)
		assert.EqualError(t, err, "cannot seek on nil file")
	})
	t.Run("nil file", func(t *testing.T) {
		file := &vfile{}
		sought, err := file.Seek(10, 10)
		assert.Equal(t, int64(0), sought)
		assert.EqualError(t, err, "cannot seek on nil file")
	})
	t.Run("non-nil file", func(t *testing.T) {
		file := &vfile{
			file: fs.NewFile("", nil, nil, nil, false),
		}
		sought, err := file.Seek(10, 10)
		assert.Equal(t, int64(0), sought)
		assert.EqualError(t, err, "bytes.Reader.Seek: invalid whence")
	})
}

func TestRead(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		var file *vfile
		buf := make([]byte, 20)
		read, err := file.Read(buf)
		assert.Equal(t, 0, read)
		assert.EqualError(t, err, "cannot read on nil file")
	})
	t.Run("nil file", func(t *testing.T) {
		file := &vfile{}
		buf := make([]byte, 20)
		read, err := file.Read(buf)
		assert.Equal(t, 0, read)
		assert.EqualError(t, err, "cannot read on nil file")
	})
	t.Run("non-nil file", func(t *testing.T) {
		file := &vfile{
			file: fs.NewFile("", nil, nil, nil, false),
		}
		buf := make([]byte, 20)
		read, err := file.Read(buf)
		assert.Equal(t, 0, read)
		assert.EqualError(t, err, "EOF")
	})
}

func TestClose(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		var file *vfile
		err := file.Close()
		assert.NoError(t, err)
	})
	t.Run("nil file", func(t *testing.T) {
		file := &vfile{}
		err := file.Close()
		assert.NoError(t, err)
	})
	t.Run("non-nil file", func(t *testing.T) {
		file := &vfile{
			file: fs.NewFile("", nil, nil, nil, true),
		}
		err := file.Close()
		assert.NoError(t, err)
	})
}

func TestReaddir(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		var file *vfile
		list, err := file.Readdir(12)
		assert.Nil(t, list)
		assert.EqualError(t, err, "cannot read directory on nil file system")
	})
	t.Run("nil filesystem", func(t *testing.T) {
		file := &vfile{}
		list, err := file.Readdir(12)
		assert.Nil(t, list)
		assert.EqualError(t, err, "cannot read directory on nil file system")
	})
	t.Run("non-nil filesystem", func(t *testing.T) {
		mfs := fs.New()
		mfs.On("ReadDir", "filename.ext").Return([]os.FileInfo{nil, nil, nil, nil}, nil)
		file := &vfile{
			fs:   mfs,
			path: "filename.ext",
		}
		list, err := file.Readdir(2)
		assert.Equal(t, []os.FileInfo{nil, nil}, list)
		assert.NoError(t, err)
	})
}

func TestStat(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		var file *vfile
		info, err := file.Stat()
		assert.Nil(t, info)
		assert.EqualError(t, err, "cannot get stats on nil file system")
	})
	t.Run("nil file system", func(t *testing.T) {
		file := &vfile{}
		info, err := file.Stat()
		assert.Nil(t, info)
		assert.EqualError(t, err, "cannot get stats on nil file system")
	})
	t.Run("non-nil file, empty path", func(t *testing.T) {
		mfs := fs.New()
		mfs.On("Stat", "").Return(nil, nil)
		file := &vfile{
			fs: mfs,
		}
		info, err := file.Stat()
		assert.Nil(t, info)
		assert.EqualError(t, err, "cannot get stats on nil file system")
	})
	t.Run("non-nil file", func(t *testing.T) {
		mfs := fs.New()
		mfs.On("Stat", "a").Return(nil, nil)
		file := &vfile{
			fs:   mfs,
			path: "a",
		}
		info, err := file.Stat()
		assert.Nil(t, info)
		assert.NoError(t, err)
	})
}

func TestGetPath(t *testing.T) {
	var f *vfile
	assert.Equal(t, "", f.getPath())
}
