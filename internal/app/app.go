package app

import "log"

type app struct {
}

// Start starts the server
func Start() {
	a := &app{}
	a.Start()
}

func (a *app) Start() {
	log.Println("hello world")
}
