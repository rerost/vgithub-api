package main

import (
	"os"

	"github.com/rerost/vgithub-api/app"
)

func main() {
	os.Exit(run())
}

func run() int {
	err := app.Run()
	if err != nil {
		return 1
	}
	return 0
}
