package e2e

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/fa93hws/run-merged-step/services"
)

func runCommand(commands []string) {
	currentDir, _ := os.Getwd()
	repoRootDir := filepath.Dir(currentDir)
	execService := &services.ExecService{}
	exitCode := execService.Run(commands[0], commands[1:], &repoRootDir)
	if exitCode != 0 {
		panic(fmt.Sprintf("command %s failed with exit code %d", commands, exitCode))
	}
}

func exists(file string) bool {
	_, err := os.Stat(file)
	return !errors.Is(err, os.ErrNotExist)
}

func getBinaryPath() string {
	if os.Getenv("BAZEL_E2E") == "true" {
		binPathInSandbox := os.Getenv("RUN_MERGED_STEP_BIN")
		if binPathInSandbox == "" {
			panic(fmt.Errorf("RUN_MERGED_STEP_BIN need to be set when run in bazel"))
		}
		return binPathInSandbox
	}

	currentDir, _ := os.Getwd()
	fixtureDir := path.Join(currentDir, "fixtures")
	outName := path.Join(fixtureDir, "run_merged_step.out")
	runCommand([]string{"go", "build", "-o", outName, "main.go"})
	return outName
}
