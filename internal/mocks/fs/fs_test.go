package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/godoc/vfs"
)

func TestNew(t *testing.T) {
	m := New()
	assert.Equal(t, &MockFS{}, m)
	m.AssertExpectations(t)
}

func TestOpen(t *testing.T) {
	m := New()
	m.On("Open", "filename.ext").Return(vfs.ReadSeekCloser(nil), nil)
	file, err := m.Open("filename.ext")
	assert.Nil(t, file)
	assert.NoError(t, err)
}

func TestLstat(t *testing.T) {
	m := New()
	m.On("Lstat", "filename.ext").Return(nil, nil)
	stat, err := m.Lstat("filename.ext")
	assert.Nil(t, stat)
	assert.NoError(t, err)
}

func TestStat(t *testing.T) {
	m := New()
	m.On("Stat", "filename.ext").Return(nil, nil)
	stat, err := m.Stat("filename.ext")
	assert.Nil(t, stat)
	assert.NoError(t, err)
}

func TestReadDir(t *testing.T) {
	m := New()
	m.On("ReadDir", "filename.ext").Return(nil, nil)
	files, err := m.ReadDir("filename.ext")
	assert.Nil(t, files)
	assert.NoError(t, err)
}

func TestRootType(t *testing.T) {
	m := New()
	m.On("RootType", "filename.ext").Return(vfs.RootType(""))
	rtype := m.RootType("filename.ext")
	assert.Empty(t, rtype)
}

func TestString(t *testing.T) {
	m := New()
	m.On("String").Return("")
	s := m.String()
	assert.Empty(t, s)
}
