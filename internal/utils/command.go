package utils

import (
	"os"
	"os/exec"
)

type (
	CommandExecutor interface {
		Run(cmd *exec.Cmd) error
	}

	execCommandExecutor struct {
	}
)

func NewCommandRunner() CommandExecutor {
	return &execCommandExecutor{}
}

func (commandRunner *execCommandExecutor) Run(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
