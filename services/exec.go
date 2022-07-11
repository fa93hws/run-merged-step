package services

import (
	"os"
	"os/exec"
)

type IExecService interface {
	Run(program string, args []string, cwd *string) int
}

type ExecService struct{}

func (e *ExecService) Run(program string, args []string, cwd *string) int {
	cmd := exec.Command(program, args...)
	if cwd != nil {
		cmd.Dir = *cwd
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		} else {
			return 1
		}
	}
	return 0
}
