package vnet

import (
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
	mfs := fs.New()
	mfs.On("Open", "filepath.ext").Return()
	v := NewDir(mfs)
	v.Open("filepath.ext")
}
