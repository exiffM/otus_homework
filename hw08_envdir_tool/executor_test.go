package main

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Run cmd", func(t *testing.T) {
		cmd := exec.Command("/bin/bash", "testdata/echo.sh", "1", "2")
		cmd.Env = append(os.Environ(), "BAR=bar", "FOO=   foo\x00with new line",
			"HELLO=\"hello\"")
		err := cmd.Run()
		expectedCode := 0
		var errExit *exec.ExitError
		if errors.As(err, &errExit) {
			expectedCode = errExit.ExitCode()
		}
		env, err := ReadDir("testdata/env")
		if err != nil {
			return
		}
		retCode := RunCmd([]string{"/bin/bash", "testdata/echo.sh", "1", "2"}, env)
		require.Equal(t, expectedCode, retCode, "actual exit code is %v", retCode)
	})
}
