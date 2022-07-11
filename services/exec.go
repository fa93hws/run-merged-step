package services

import (
	"fmt"
	"os"
	"os/exec"
)

type IExecService interface {
	Run(program string, args []string, cwd *string, env *map[string]string) int
}

type ExecService struct{}

func (e *ExecService) Run(program string, args []string, cwd *string, env *map[string]string) int {
	cmd := exec.Command(program, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if cwd != nil {
		cmd.Dir = *cwd
	}
	if env != nil {
		for key, val := range *env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, val))
		}
	}
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
