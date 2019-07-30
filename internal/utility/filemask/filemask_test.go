package filemask

import (
	"runtime/debug"
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/fs"
	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	t.Run("no path change", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Failf(t, "caught panic", "%v\n%s", r, debug.Stack())
			}
		}()
		mfs := fs.New()
		mfs.On("Open", "myfile.ext").Return(new(fs.MockFile))
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
		mfs.On("Open", "inner/path/to/something/myfile.ext").Return(new(fs.MockFile))
		newFs := Wrap(mfs, "inner/path/to/something/")
		file, err := newFs.Open("myfile.ext")
		assert.Equal(t, new(fs.MockFile), file)
		assert.NoError(t, err)
	})
}
