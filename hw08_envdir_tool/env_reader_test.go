package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Все сценарии", func(t *testing.T) {
		expectEnv := Environment{
			"BAR":   EnvValue{Value: "bar"},
			"EMPTY": EnvValue{Value: ""},
			"FOO":   EnvValue{Value: "   foo\nwith new line"},
			"HELLO": EnvValue{Value: `"hello"`},
			"UNSET": EnvValue{NeedRemove: true},
		}
		env, err := ReadDir("testdata/env")
		require.NoError(t, err)
		require.Equal(t, env, expectEnv)
	})

	t.Run("Пустая директория", func(t *testing.T) {
		emptyDir := "testdata/empty"
		if err := os.Mkdir(emptyDir, 0o755); err != nil {
			t.Fatalf("Error creating directory: %v", err)
		}
		defer func() {
			if err := os.Remove(emptyDir); err != nil {
				t.Fatalf("Error removing directory: %v", err)
			}
		}()

		expectEnv := Environment{}

		env, err := ReadDir(emptyDir)
		require.NoError(t, err)
		require.Equal(t, env, expectEnv)
	})

	t.Run("Директория недоступная для чтения", func(t *testing.T) {
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

		env, err := ReadDir(notReadableDir)
		require.Error(t, err)
		require.Equal(t, "directory is not readable: testdata/not_readable", err.Error())
		require.Nil(t, env)
	})

	t.Run("Директория не существует", func(t *testing.T) {
		dirNotExist := "testdata/not_exist"

		env, err := ReadDir(dirNotExist)
		require.Error(t, err)
		require.Equal(t, "directory does not exist: testdata/not_exist", err.Error())
		require.Nil(t, env)
	})

	t.Run("Это не директория", func(t *testing.T) {
		fileAsDirectory := "testdata/not_exist"
		// Создать файл с именем директории
		if _, err := os.Create(fileAsDirectory); err != nil {
			t.Fatalf("Error creating file: %v", err)
		}
		defer func() {
			if err := os.Remove(fileAsDirectory); err != nil {
				t.Fatalf("Error removing file: %v", err)
			}
		}()

		env, err := ReadDir(fileAsDirectory)
		require.Error(t, err)
		require.Equal(t, "not a directory: testdata/not_exist", err.Error())
		require.Nil(t, env)
	})
}
