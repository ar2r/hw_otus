package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {

	tests := []struct {
		name     string
		envDir   string
		envFile  string
		args     []string
		expected string
	}{
		{
			name:     "BAR - несколько строк",
			envFile:  "testdata/env/BAR",
			args:     []string{"arg1=1", "arg2=2"},
			expected: "HELLO is ()\nBAR is (bar)\nFOO is ()\nUNSET is (unset)\nADDED is ()\nEMPTY is ()\narguments are arg1=1 arg2=2\n",
		},
		{
			name:     "EMPTY - несколько пустых строк с пробелами",
			envFile:  "testdata/env/EMPTY",
			args:     []string{"arg1=1", "arg2=2"},
			expected: "HELLO is ()\nBAR is ()\nFOO is ()\nUNSET is (unset)\nADDED is ()\nEMPTY is ()\narguments are arg1=1 arg2=2\n",
		},
		{
			name:     "FOO - NUL символ",
			envFile:  "testdata/env/FOO",
			args:     []string{"arg1=1", "arg2=2"},
			expected: "HELLO is ()\nBAR is ()\nFOO is (   foo\nwith new line)\nUNSET is (unset)\nADDED is ()\nEMPTY is ()\narguments are arg1=1 arg2=2\n",
		},
		{
			name:     "HELLO - просто одна строка",
			envFile:  "testdata/env/HELLO",
			args:     []string{"arg1=1", "arg2=2"},
			expected: "HELLO is (\"hello\")\nBAR is ()\nFOO is ()\nUNSET is (unset)\nADDED is ()\nEMPTY is ()\narguments are arg1=1 arg2=2\n",
		},
		{
			name:     "UNSET - Удаление ENV",
			envFile:  "testdata/env/UNSET",
			args:     []string{"arg1=1", "arg2=2"},
			expected: "HELLO is ()\nBAR is ()\nFOO is ()\nUNSET is ()\nADDED is ()\nEMPTY is ()\narguments are arg1=1 arg2=2\n",
		},
		{
			name:     "TestHappyPath",
			envDir:   "testdata/env",
			args:     []string{"arg1=1", "arg2=2"},
			expected: "HELLO is (\"hello\")\nBAR is (bar)\nFOO is (   foo\nwith new line)\nUNSET is ()\nADDED is ()\nEMPTY is ()\narguments are arg1=1 arg2=2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохранить оригинальный stdout и восстановить его после теста
			originalStdout := os.Stdout
			defer func() { os.Stdout = originalStdout }()
			// Создать pipe для перехвата stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			if tt.envFile != "" {
				// Создать временную директорию и скопировать туда tt.envFile
				envDir := t.TempDir()
				// Скпировать файл в директорию
				err := copyFile(tt.envFile, envDir)

				require.NoError(t, err)
				tt.envDir = envDir
				// Удалить временную директорию после теста
				defer os.RemoveAll(envDir)
			}

			// Тест
			os.Setenv("UNSET", "unset")
			err := runCommand("testdata/echo.sh", getEnvironment(tt.envDir), tt.args)
			require.NoError(t, err)

			// Закрыть writer и прочитать stdout
			w.Close()
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			require.Equal(t, tt.expected, output)
		})
	}
}

func copyFile(src, dstDir string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstPath := filepath.Join(dstDir, filepath.Base(src))
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	err = dstFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
