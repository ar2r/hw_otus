package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Создать пустую директори
	emptyDir := "testdata/empty"
	if err := os.Mkdir(emptyDir, 0o755); err != nil {
		t.Fatalf("Error creating directory: %v", err)
	}
	defer func() {
		if err := os.Remove(emptyDir); err != nil {
			t.Fatalf("Error removing directory: %v", err)
		}
	}()

	// Создать директорию недостпуную для чтения
	notReadableDir := "testdata/not_readable"
	if err := os.Mkdir(notReadableDir, 0o000); err != nil {
		t.Fatalf("Error creating directory: %v", err)
	}
	defer func() {
		if err := os.Chmod(notReadableDir, 0o755); err != nil {
			t.Fatalf("Error changing directory permissions: %v", err)
		}
		if err := os.Remove(notReadableDir); err != nil {
			t.Fatalf("Error removing directory: %v", err)
		}
	}()

	tests := []struct {
		name string
		dir  string
		err  string
	}{
		{
			name: "Happy path",
			dir:  "testdata/env",
			err:  "",
		},
		{
			name: "Пустая директория",
			dir:  emptyDir,
			err:  "",
		},
		{
			name: "Директория не существует",
			dir:  "testdata/not_exist",
			err:  "directory does not exist: testdata/not_exist",
		},
		{
			name: "Не директория",
			dir:  "testdata/echo.sh",
			err:  "not a directory: testdata/echo.sh",
		},
		{
			name: "Директория не доступна для чтения",
			dir:  notReadableDir,
			err:  "directory is not readable: testdata/not_readable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := ReadDir(tt.dir)

			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				if err.Error() != tt.err {
					t.Errorf("ReadDir(%s) failed: got [%v], want [%v]", tt.dir, err, tt.err)
				}
			}

			fmt.Println("Environment:")
			for k, v := range env {
				fmt.Printf("%s: %v\n", k, v)
			}
		})
	}
}
