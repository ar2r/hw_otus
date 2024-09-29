package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name     string
		cmd      []string
		env      Environment
		expected int
		stdout   string
	}{
		{
			name:     "Простой echo",
			cmd:      []string{"echo", "hello"},
			env:      Environment{},
			expected: 0,
			stdout:   "hello\n",
		},
		{
			name:     "Использование аргумента",
			cmd:      []string{"sh", "-c", "echo hello"},
			env:      Environment{},
			expected: 0,
			stdout:   "hello\n",
		},
		{
			name:     "Использование ENV",
			cmd:      []string{"sh", "-c", "echo hello $NAME!"},
			env:      Environment{"NAME": {"world", false}},
			expected: 0,
			stdout:   "hello world!\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			r, w, _ := os.Pipe()
			oldStdout := os.Stdout
			defer func() { os.Stdout = oldStdout }()
			os.Stdout = w

			got := RunCmd(tt.cmd, tt.env)
			w.Close()
			io.Copy(&stdout, r)

			if got != tt.expected {
				t.Errorf("RunCmd() = %v, want %v", got, tt.expected)
			}
			if stdout.String() != tt.stdout {
				t.Errorf("stdout = %v, want %v", stdout.String(), tt.stdout)
			}
		})
	}
}
