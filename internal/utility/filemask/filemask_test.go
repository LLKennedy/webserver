package filemask

import (
	"runtime/debug"
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/fs"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/godoc/vfs"
)

func TestWrap(t *testing.T) {
	t.Run("root path", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Failf(t, "caught panic", "%v\n%s", r, debug.Stack())
			}
		}()
		mfs := fs.New()
		mfs.On("Open", "/index.html").Return(new(fs.MockFile), nil)
		newFs := Wrap(mfs, "")
		file, err := newFs.Open("/")
		assert.Equal(t, new(fs.MockFile), file)
		assert.NoError(t, err)
	})
	t.Run("no path change", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Failf(t, "caught panic", "%v\n%s", r, debug.Stack())
			}
		}()
		mfs := fs.New()
		mfs.On("Open", "myfile.ext").Return(new(fs.MockFile), nil)
		newFs := Wrap(mfs, "")
		file, err := newFs.Open("myfile.ext")
		assert.Equal(t, new(fs.MockFile), file)
		assert.NoError(t, err)
	})
	t.Run("change path", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Failf(t, "caught panic", "%v\n%s", r, debug.Stack())
			}
		}()
		mfs := fs.New()
		mfs.On("Open", "inner/path/to/something/myfile.ext").Return(new(fs.MockFile), nil)
		newFs := Wrap(mfs, "inner/path/to/something/")
		file, err := newFs.Open("myfile.ext")
		assert.Equal(t, new(fs.MockFile), file)
		assert.NoError(t, err)
	})
}

func TestLstat(t *testing.T) {
	var f *filemask
	osInfo, osErr := vfs.OS(".").Lstat("")
	info, err := f.Lstat("")
	assert.Equal(t, osInfo, info)
	if osErr == nil {
		assert.NoError(t, err)
	} else {
		assert.EqualError(t, err, osErr.Error())
	}
}

func TestStat(t *testing.T) {
	var f *filemask
	osInfo, osErr := vfs.OS(".").Stat("")
	info, err := f.Stat("")
	assert.Equal(t, osInfo, info)
	if osErr == nil {
		assert.NoError(t, err)
	} else {
		assert.EqualError(t, err, osErr.Error())
	}
}

func TestReadDir(t *testing.T) {
	var f *filemask
	osInfo, osErr := vfs.OS(".").ReadDir("")
	info, err := f.ReadDir("")
	assert.Equal(t, osInfo, info)
	if osErr == nil {
		assert.NoError(t, err)
	} else {
		assert.EqualError(t, err, osErr.Error())
	}
}

func TestRootType(t *testing.T) {
	var f *filemask
	osInfo := vfs.OS(".").RootType("")
	info := f.RootType("")
	assert.Equal(t, osInfo, info)
}

func TestString(t *testing.T) {
	var f *filemask
	osStr := vfs.OS(".").String()
	str := f.String()
	assert.Equal(t, osStr, str)
}
