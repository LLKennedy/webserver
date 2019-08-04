package app

import (
	"fmt"
	"testing"

	"github.com/LLKennedy/webserver/internal/mocks/fs"
	"github.com/LLKennedy/webserver/internal/mocks/mocknetwork"
	"github.com/LLKennedy/webserver/internal/network"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	mNet := new(mocknetwork.Layer)
	mNet2 := new(mocknetwork.Layer)
	secure = mNet
	insecure = mNet2
	mNet.On("ListenAndServe", "localhost:80", network.NewHTTPServer(logger, fileSystem, secure, insecure)).Return(fmt.Errorf("cannot start"))
	fileSystem = fs.New(fs.NewFile("build/index.html", nil, fmt.Errorf("the system cannot find the path specified"), nil, false))
	err := Run()
	assert.EqualError(t, err, "could not read script hash: could not open index file: the system cannot find the path specified")
}
