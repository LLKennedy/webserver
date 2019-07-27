package vnet

import (
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/fs"
)

func TestGetFs(t *testing.T) {
	v := &VDir{
		fs: fs.New(),
	}
	v.
}
