package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	defaults := Options{
		Address:      "localhost",
		Port:         443,
		InsecurePort: 80,
		KeyFile:      defaultKeyFile,
		CertFile:     defaultCertFile,
	}
	loaded, err := Load(nil, nil)
	assert.Equal(t, defaults, loaded)
	assert.NoError(t, err)
}
