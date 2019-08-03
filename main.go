package main

import (
	"fmt"
	"os"

	"github.com/LLKennedy/webserver/internal/app"
)

var (
	exitProgram = os.Exit
	runApp      = app.Run
)

func main() {
	if err := runApp(); err != nil {
		fmt.Printf("error starting application: %v", err)
		exitProgram(1)
	}
	exitProgram(0)
}
