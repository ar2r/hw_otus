package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	if err := isDirectoryCompatible(dir); err != nil {
		return nil, err
	}

	// Прочитать содержимое директории
	files, err := getFilePathSlice(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filePath := fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), file.Name())
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %s, error: %s", filePath, err)
		}

		if len(fileContent) == 0 {
			env[file.Name()] = EnvValue{Value: "", NeedRemove: true}
			continue
		}

		fileContent = findUniversalNewline(fileContent)
		fileContent = bytes.Replace(fileContent, []byte{0}, []byte{'\n'}, -1)
		fileContent = []byte(strings.TrimRight(string(fileContent), " \t\n\r"))
		env[file.Name()] = EnvValue{Value: string(fileContent), NeedRemove: false}
	}

	return env, nil
}

func getFilePathSlice(dir string) ([]os.DirEntry, error) {
	compatibleFiles := make([]os.DirEntry, 0)

	originalFiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %s", dir)
	}

	for _, file := range originalFiles {
		// Пропустить директории
		if file.IsDir() {
			continue
		}

		// Пропустить файлы начинающиеся с точки (скрытые или системные файлы)
		if file.Name()[0] == '.' {
			continue
		}

		// Проверить доступность файла на чтение
		if _, err := os.Open(fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), file.Name())); err != nil {
			return nil, fmt.Errorf("error reading env file: %s", file.Name())
		}

		compatibleFiles = append(compatibleFiles, file)
	}

	return compatibleFiles, nil
}

func isDirectoryCompatible(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	if fileInfo, err := os.Stat(dir); err != nil || !fileInfo.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}

	if _, err := os.Open(dir); err != nil {
		return fmt.Errorf("directory is not readable: %s", dir)
	}

	return nil
}

func findUniversalNewline(fileContent []byte) []byte {
	for i := 0; i < len(fileContent); i++ {
		if fileContent[i] == '\n' {
			return fileContent[:i]
		}
		if fileContent[i] == '\r' {
			if i+1 < len(fileContent) && fileContent[i+1] == '\n' {
				return fileContent[:i]
			}
			return fileContent[:i]
		}
	}

	return fileContent
}
