package app

import (
	"log"
	"os"

	"github.com/LLKennedy/webserver/internal/network"
	"github.com/LLKennedy/webserver/internal/utility/logs"
	"golang.org/x/tools/godoc/vfs"
)

var (
	logger     logs.Logger    = log.New(os.Stdout, "webserver ", log.Flags())
	fileSystem vfs.FileSystem = vfs.OS(".")
	net        network.Layer  = network.HTTP{}
)

// App is the application
type App struct {
	HTTPServer Server
}

// Run creates, configures and starts the app
func Run() error {
	a := New(logger, fileSystem, net)
	return a.Start()
}

// New creates a new app struct
func New(lg logs.Logger, fs vfs.FileSystem, net network.Layer) *App {
	a := &App{
		HTTPServer: network.NewHTTPServer(lg, fs, net),
	}
	return a
}

// Start starts the app
func (a *App) Start() error {
	return a.HTTPServer.Start()
}
