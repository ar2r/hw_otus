package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

var exitCode int

func fakeExit(code int) {
	exitCode = code
}

var osExit = os.Exit

func TestRun(t *testing.T) {
	// Сохранить оригинальную функцию os.Exit и восстановить ее после теста
	originalExit := osExit
	defer func() { osExit = originalExit }()
	osExit = fakeExit

	// Сохранить оригинальный stdout и восстановить его после теста
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()

	// Создать pipe для перехвата stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	t.Run("TestHappyPath", func(t *testing.T) {
		// mock os.Args
		os.Args = []string{"main", "testdata/env", "/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"}

		main()

		// Закрыть writer и прочитать stdout
		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		expected := "HELLO is (\"hello\")\nBAR is (bar)\nFOO is (   foo\nwith new line)\nUNSET is ()\nADDED is ()\nEMPTY is ()\narguments are arg1=1 arg2=2\n"

		require.Equal(t, expected, output)
		require.Equal(t, 0, exitCode)
	})
}
