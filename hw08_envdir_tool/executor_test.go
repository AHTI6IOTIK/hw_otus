package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmdQuit(t *testing.T) {
	args := make([]string, 0, 1)
	envs := make(Environment)
	actual := RunCmd(args, envs)

	assert.Equal(t, InvalidCmdArgs, actual)
}
