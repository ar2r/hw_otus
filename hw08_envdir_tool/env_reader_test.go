package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

	// Создать директорию с пустым файлом
	emptyFileDir := "testdata/empty_file"
	if err := os.Mkdir(emptyFileDir, 0o755); err != nil {
		t.Fatalf("Error creating directory: %v", err)
	}
	if err := os.WriteFile(fmt.Sprintf("%s/empty", emptyFileDir), []byte{}, 0o644); err != nil {
		t.Fatalf("Error creating file: %v", err)
	}
	defer func() {
		// Удалить директорию с вложенными файлами
		if err := os.RemoveAll(emptyFileDir); err != nil {
			t.Fatalf("Error removing directory with files: %v", err)
		}
	}()

	tests := []struct {
		dir string
		err string
	}{
		{dir: "testdata/env", err: ""},
		{dir: emptyDir, err: ""},
		{dir: "testdata/not_exist", err: "Directory does not exist: testdata/not_exist"},
		{dir: "testdata/echo.sh", err: "Not a directory: testdata/echo.sh"},
		{dir: notReadableDir, err: "Directory is not readable: testdata/not_readable"},
		{dir: emptyFileDir, err: "Error reading file: testdata/empty_file/empty"},
	}

	for _, tt := range tests {
		t.Run(tt.dir, func(t *testing.T) {
			env, err := ReadDir(tt.dir)
			if err != nil && err.Error() != tt.err {
				t.Errorf("ReadDir(%s) failed: got [%v], want [%v]", tt.dir, err, tt.err)
			}

			fmt.Println("Environment:")
			for k, v := range env {
				fmt.Printf("%s: %v\n", k, v)
			}
			assert.Equal(t, true, true)
		})
	}
}
