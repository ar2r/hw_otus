package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrNotEnoughArguments = errors.New("not enough arguments")
	ErrReadingEnvDir      = errors.New("error reading env dir")
	ErrGettingAbsPath     = errors.New("error getting absolute path")
	ErrGettingEnvDirPath  = errors.New("error getting env dir path")
)

func main() {
	if len(os.Args) < 3 {
		printExit(ErrNotEnoughArguments, nil)
	}
	cmd := append([]string{getExecutableFilePath()}, getOsArgs()...)
	exitCode := RunCmd(cmd, getEnvironment(getEnvDirPath()))
	os.Exit(exitCode)
}

func getOsArgs() []string {
	var args []string

	for i := 3; i < len(os.Args); i++ {
		args = append(args, os.Args[i])
	}
	return args
}

func getEnvironment(envDirPath string) Environment {
	myEnv, err := ReadDir(envDirPath)
	if err != nil {
		printExit(ErrReadingEnvDir, err)
	}
	return myEnv
}

func getExecutableFilePath() string {
	executableFilePath := os.Args[2]
	executableFilePath, err := filepath.Abs(executableFilePath)
	if err != nil {
		printExit(ErrGettingAbsPath, err)
	}
	return executableFilePath
}

func getEnvDirPath() string {
	envDirPath := os.Args[1]
	envDirPath, err := filepath.Abs(envDirPath)
	if err != nil {
		printExit(ErrGettingEnvDirPath, err)
	}
	return envDirPath
}

func printExit(msg error, err error) {
	fmt.Println(msg, err)
	os.Exit(1)
}
