package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	ErrNotEnoughArguments = errors.New("not enough arguments")
	ErrReadingEnvDir      = errors.New("error reading env dir")
	ErrRunCommand         = errors.New("error running command")
	ErrGettingAbsPath     = errors.New("error getting absolute path")
	ErrGettingEnvDirPath  = errors.New("error getting env dir path")
)

func main() {
	if len(os.Args) < 3 {
		printExit(ErrNotEnoughArguments, nil)
	}

	envDirPath := getEnvDirPath()
	executableFilePath := getExecutableFilePath()
	myEnv := getEnvironment(envDirPath)
	args := getArgs()

	// Запустить команду
	cmd := exec.Command(executableFilePath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Pass env which was set before via export
	for _, v := range os.Environ() {
		cmd.Env = append(cmd.Env, v)
	}

	// Передать myEnv в исполняемый файл
	for k, v := range myEnv {
		if v.NeedRemove {
			cmd.Env = append(cmd.Env, k+"=")
			continue
		}
		cmd.Env = append(cmd.Env, k+"="+v.Value)
	}

	// Запустить команду
	errCmd := cmd.Run()
	if errCmd != nil {
		if exitError, ok := errCmd.(*exec.ExitError); ok {
			// Вернуть exit code выполненной команды
			os.Exit(exitError.ExitCode())
		} else {
			printExit(ErrRunCommand, errCmd)
		}
	}
}

func getArgs() []string {
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
