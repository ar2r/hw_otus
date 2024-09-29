package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrNotEnoughArguments = errors.New("not enough arguments")
	ErrReadingEnvDir      = errors.New("error reading env dir")
)

func main() {
	if len(os.Args) < 2 {
		printExit(ErrNotEnoughArguments, nil)
	}
	exitCode := RunCmd(os.Args[2:], getEnvironment(os.Args[1]))
	os.Exit(exitCode)
}

func getEnvironment(envDirPath string) Environment {
	myEnv, err := ReadDir(envDirPath)
	if err != nil {
		printExit(ErrReadingEnvDir, err)
	}
	return myEnv
}

func printExit(msg error, err error) {
	fmt.Println(msg, err)
	os.Exit(1)
}
