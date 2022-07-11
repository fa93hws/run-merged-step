package e2e

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func runCommand(commands []string) {
	currentDir, _ := os.Getwd()
	repoRootDir := filepath.Dir(currentDir)
	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Dir = repoRootDir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func getStatusFile(jobId string) string {
	return path.Join(os.TempDir(), jobId, "merged_step_status.json")
}

func exists(file string) bool {
	_, err := os.Stat(file)
	return !errors.Is(err, os.ErrNotExist)
}

func readContent(file string) string {
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(content)
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
