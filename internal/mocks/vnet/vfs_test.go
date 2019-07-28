package vnet

import (
	"fmt"
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

}
