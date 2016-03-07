package main

import (
	"os"
)

const Name = "mantl-bootstrap"
const Version = "0.1.0"

func main() {
	root := initCommand(Name, Version)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
