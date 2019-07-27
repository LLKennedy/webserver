package vnet

import (
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/fs"
	"github.com/stretchr/testify/assert"
)

func TestNewDir(t *testing.T) {
	nfs := fs.New()
	v := NewDir(nfs)
	expected := new(VDir)
	expected.fs = nfs
	assert.Equal(t, expected, v)
}
