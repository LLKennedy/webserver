package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPathNode(t *testing.T) {
	t.Run("short", func(t *testing.T) {
		inString := "a"
		head, tail := getPathNode(inString)
		assert.Equal(t, "a", head)
		assert.Equal(t, "/", tail)
	})
	t.Run("long", func(t *testing.T) {
		inString := "a/b/c"
		head, tail := getPathNode(inString)
		assert.Equal(t, "a", head)
		assert.Equal(t, "/b/c", tail)
	})
}
