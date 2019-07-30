package app

import (
	"fmt"
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/mocknetwork"
	"github.com/LLKennedy/webserver/internal/network"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	mNet := new(mocknetwork.Layer)
	net = mNet
	mNet.On("ListenAndServe", "localhost:80", network.NewHTTPServer(logger, fileSystem, net)).Return(fmt.Errorf("cannot start"))
	err := Run()
	assert.EqualError(t, err, "http server closed unexpectedly: cannot start")
}
