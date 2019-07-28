package app

import (
	"github.com/LLKennedy/webserver/internal/network"
	"golang.org/x/tools/godoc/vfs"
)

var (
	fileSystem vfs.FileSystem = vfs.OS(".")
	net        network.Layer  = network.HTTP{}
)

// App is the application
type App struct {
	HTTPServer Server
}

// Server starts and only exits if there's an error
type Server interface {
	Start() error
}

// Run creates, configures and starts the app
func Run() error {
	a := New(fileSystem, net)
	return a.Start()
}

// New creates a new app struct
func New(fs vfs.FileSystem, net network.Layer) *App {
	a := &App{
		HTTPServer: network.NewHTTPServer(fs, net),
	}
	return a
}

// Start starts the app
func (a *App) Start() error {
	return a.HTTPServer.Start()
}
