package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultOptions(t *testing.T) {
	defaults := Options{
		Address:      "localhost",
		Port:         443,
		InsecurePort: 80,
		KeyFile:      defaultKeyFile,
		CertFile:     defaultCertFile,
	}
	assert.Equal(t, defaults, DefaultOptions())
}
