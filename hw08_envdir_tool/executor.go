package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	execCmd := exec.Command(cmd[0], cmd[1:]...) //nolint: gosec
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	// Передать текущие переменные окружения
	execCmd.Env = append(execCmd.Env, os.Environ()...)

	// Передать переменные окружения из файлов
	for k, v := range env {
		if v.NeedRemove {
			execCmd.Env = append(execCmd.Env, k+"=")
			continue
		}
		execCmd.Env = append(execCmd.Env, k+"="+v.Value)
	}

	// Запустить команду
	errCmd := execCmd.Run()
	if errCmd != nil {
		var exitError *exec.ExitError
		if errors.As(errCmd, &exitError) {
			// Вернуть exit code выполненной команды
			return exitError.ExitCode()
		}
	}

	return 0
}
