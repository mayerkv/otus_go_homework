package main

import (
	"errors"
	"log"
	"os"
)

var ErrInvalidArgs = errors.New("invalid arguments")

func main() {
	if len(os.Args) < 3 {
		log.Fatal(ErrInvalidArgs)
	}

	path := os.Args[1]
	command := os.Args[2:]

	environment, err := ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(RunCmd(command, environment))
}
