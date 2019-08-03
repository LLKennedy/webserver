package config

import (
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/fs"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	mfs := fs.New(fs.NewFile("%APPDATA%/webserver/config.json", []byte("{}"), nil, nil, true))
	defaults := Options{
		Address:      "localhost",
		Port:         443,
		InsecurePort: 80,
		KeyFile:      defaultKeyFile,
		CertFile:     defaultCertFile,
	}
	loaded, err := Load(mfs, nil)
	assert.Equal(t, defaults, loaded)
	assert.NoError(t, err)
}
