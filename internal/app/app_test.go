package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	err := Run()
	assert.EqualError(t, err, "http server closed unexpectedly: listen tcp: address localhost: missing port in address")
}
