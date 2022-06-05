package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	name := cmd[0]
	args := cmd[1:]

	command := exec.Command(name, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, v.Value)
		}
	}

	command.Env = os.Environ()

	if err := command.Run(); err != nil {
		return 1
	}

	return 0
}
