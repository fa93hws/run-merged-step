package services

import (
	"os"
	"os/exec"
)

type IExecService interface {
	Run(program string, args []string) int
}

type ExecService struct {
	cwd *string
}

func NewExecService(cwd *string) *ExecService {
	return &ExecService{cwd}
}

func (e *ExecService) Run(program string, args []string) int {
	cmd := exec.Command(program, args...)
	if e.cwd != nil {
		cmd.Dir = *e.cwd
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
