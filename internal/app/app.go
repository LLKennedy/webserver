package app

import (
	"log"

	"golang.org/x/tools/godoc/vfs"
)

var (
	fileSystem = vfs.OS(".")
)

// App is the application
type App struct {
	fs vfs.FileSystem
}

// Run creates, configures and starts the app
func Run() {
	a := New(fileSystem)
	a.Start()
}

// New creates a new app struct
func New(fs vfs.FileSystem) *App {
	a := &App{
		fs: fs,
	}
	return a
}

// Start starts the app
func (a *App) Start() {
	log.Println("hello world")
}
