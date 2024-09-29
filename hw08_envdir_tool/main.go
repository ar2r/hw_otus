package main

import (
	"errors"
	"log"
	"os"
)

var (
	ErrNotEnoughArguments = errors.New("not enough arguments")
	ErrReadingEnvDir      = errors.New("error reading env dir")
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(ErrNotEnoughArguments)
	}
	exitCode := RunCmd(os.Args[2:], getEnvironment(os.Args[1]))
	os.Exit(exitCode)
}

func getEnvironment(envDirPath string) Environment {
	myEnv, err := ReadDir(envDirPath)
	if err != nil {
		log.Fatal(ErrReadingEnvDir, err)
	}
	return myEnv
}
