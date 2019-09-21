package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultOptions(t *testing.T) {
	defaults := Options{
		Address:       "localhost",
		Port:          443,
		InsecurePort:  80,
		KeyFile:       defaultKeyFile,
		CertFile:      defaultCertFile,
		StaticContent: "build",
	}
	assert.Equal(t, defaults, DefaultOptions())
}

func TestGetAddress(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		o := Options{}
		assert.Equal(t, DefaultOptions().Address, o.GetAddress())
	})
	t.Run("non-empty options", func(t *testing.T) {
		o := Options{}
		o.Address = "something"
		assert.Equal(t, "something", o.GetAddress())
	})
}

func TestGetPort(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		o := Options{}
		assert.Equal(t, DefaultOptions().Port, o.GetPort())
	})
	t.Run("non-empty options", func(t *testing.T) {
		o := Options{}
		o.Port = 12
		assert.Equal(t, uint16(12), o.GetPort())
	})
}

func TestGetInsecurePort(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		o := Options{}
		assert.Equal(t, DefaultOptions().InsecurePort, o.GetInsecurePort())
	})
	t.Run("non-empty options", func(t *testing.T) {
		o := Options{}
		o.InsecurePort = 12
		assert.Equal(t, uint16(12), o.GetInsecurePort())
	})
}

func TestGetKeyFile(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		o := Options{}
		assert.Equal(t, DefaultOptions().KeyFile, o.GetKeyFile())
	})
	t.Run("non-empty options", func(t *testing.T) {
		o := Options{}
		o.KeyFile = "something"
		assert.Equal(t, "something", o.GetKeyFile())
	})
}

func TestGetCertFile(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		o := Options{}
		assert.Equal(t, DefaultOptions().CertFile, o.GetCertFile())
	})
	t.Run("non-empty options", func(t *testing.T) {
		o := Options{}
		o.CertFile = "something"
		assert.Equal(t, "something", o.GetCertFile())
	})
}

func TestGetStaticContent(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		o := Options{}
		assert.Equal(t, DefaultOptions().StaticContent, o.GetStaticContent())
	})
	t.Run("non-empty options", func(t *testing.T) {
		o := Options{}
		o.StaticContent = "something"
		assert.Equal(t, "something", o.GetStaticContent())
	})
}
