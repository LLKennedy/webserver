package config

import (
	"testing"

	"github.com/LLKennedy/goconfig"
	"github.com/LLKennedy/webserver/internal/mocks/fs"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	mfs := fs.New(fs.NewFile(goconfig.DefaultConfigLocation("webserver"), []byte("{}"), nil, nil, true))
	defaults := Options{
		Address:       "localhost",
		Port:          443,
		InsecurePort:  80,
		KeyFile:       defaultKeyFile,
		CertFile:      defaultCertFile,
		StaticContent: "build",
	}
	loaded, err := Load(mfs, nil)
	assert.Equal(t, defaults, loaded)
	assert.NoError(t, err)
}
