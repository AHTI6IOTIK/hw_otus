package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

const (
	InvalidCmdArgs    int = 2
	UnknownExecuteErr int = 3
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return InvalidCmdArgs
	}

	//nolint:gosec
	command := exec.Command(cmd[0], cmd[1:]...)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	for key, vr := range env {
		err := os.Unsetenv(key)
		if err != nil {
			log.Println(err)
		}

		if vr.NeedRemove {
			continue
		}

		err = os.Setenv(key, vr.Value)
		if err != nil {
			log.Println(err)
		}
	}

	var execErr *exec.ExitError
	err := command.Run()
	if errors.As(err, &execErr) {
		return execErr.ProcessState.ExitCode()
	} else if err != nil {
		log.Println(err)
		return UnknownExecuteErr
	}

	return
}
