package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	for key, vr := range env {
		if vr.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {

			}
			continue
		}

		err := os.Setenv(key, vr.Value)
		if err != nil {

		}
	}

	err := command.Run()
	if er, ok := err.(*exec.ExitError); ok {
		return er.ProcessState.ExitCode()
	} else if err != nil {
		log.Println(err)
		return 1
	}

	return
}
