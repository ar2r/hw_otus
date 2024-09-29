package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var reEnvName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

type envFile struct {
	name string
	path string
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := Environment{}

	if err := isDirectoryCompatible(dir); err != nil {
		return nil, err
	}

	// Прочитать содержимое директории
	files, err := getFilePathSlice(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileContent, err := os.ReadFile(file.path)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %s, error: %w", file.path, err)
		}

		if len(fileContent) == 0 {
			env[file.name] = EnvValue{Value: "", NeedRemove: true}
			continue
		}

		fileContent = cleanupContent(fileContent)
		env[file.name] = EnvValue{Value: string(fileContent), NeedRemove: false}
	}

	return env, nil
}

func cleanupContent(fileContent []byte) []byte {
	fileContent = truncateRightNewLine(fileContent)
	fileContent = bytes.ReplaceAll(fileContent, []byte{0}, []byte{'\n'})
	fileContent = []byte(strings.TrimRight(string(fileContent), " \t\n\r"))
	return fileContent
}

func getFilePathSlice(dir string) ([]envFile, error) {
	compatibleFiles := make([]envFile, 0)

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

		if !reEnvName.MatchString(file.Name()) {
			return nil, fmt.Errorf("invalid env file name: %s", file.Name())
		}

		envFile := envFile{
			name: file.Name(),
			path: filepath.Join(dir, file.Name()),
		}
		compatibleFiles = append(compatibleFiles, envFile)
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

func truncateRightNewLine(fileContent []byte) []byte {
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
