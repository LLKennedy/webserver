package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	runCount := 0
	runApp = func() error {
		runCount++
		return fmt.Errorf("some error")
	}
	exitCode := -1
	exitProgram = func(in int) {
		exitCode = in
	}
	main()
	assert.Equal(t, 1, runCount)
	assert.Equal(t, 0, exitCode)
	runApp = func() error {
		runCount++
		return nil
	}
	main()
	assert.Equal(t, 2, runCount)
	assert.Equal(t, 0, exitCode)
}
