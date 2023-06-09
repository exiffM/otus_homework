package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	proc := exec.Command(cmd[0], cmd[1:]...) // #nosec G204

	for key, val := range env {
		if !val.NeedRemove {
			os.Setenv(key, val.Value)
		} else {
			os.Unsetenv(key)
		}
	}
	proc.Env = os.Environ()

	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	err := proc.Run()

	var errExit *exec.ExitError
	if errors.As(err, &errExit) {
		returnCode = errExit.ExitCode()
	}
	return
}
