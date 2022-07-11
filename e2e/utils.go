package e2e

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/fa93hws/run-merged-step/services"
)

func runCommand(commands []string, env *map[string]string) int {
	currentDir, _ := os.Getwd()
	repoRootDir := filepath.Dir(currentDir)
	execService := &services.ExecService{}
	return execService.Run(commands[0], commands[1:], &repoRootDir, env)
}

func exists(file string) bool {
	_, err := os.Stat(file)
	return !errors.Is(err, os.ErrNotExist)
}

func readFileContent(file string) string {
	bytes, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(bytes)
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
	exitCode := runCommand([]string{"go", "build", "-o", outName, "main.go"}, nil)
	if exitCode != 0 {
		panic(fmt.Errorf("failed to build run_merged_step, ExitCode=%d", exitCode))
	}
	return outName
}
