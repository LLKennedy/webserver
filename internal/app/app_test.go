package app

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	Start()
	log.SetOutput(os.Stdout)
	output := buf.Bytes()
	expected := "hello world\n"
	assert.Equal(t, output[len(output)-len(expected):], []byte(expected))
}
