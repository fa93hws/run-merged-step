package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fa93hws/run-merged-step/services"
	"github.com/spf13/cobra"
)

type RunParams struct {
	label          string
	key            string
	jobId          string
	autoRevertable bool
	commands       []string
}

var (
	label          string
	key            string
	autoRevertable bool
	runCmd         = &cobra.Command{
		Use:   "run",
		Short: "Run the CI step as a command",
		Run: func(cmd *cobra.Command, args []string) {
			osFs := &services.OsFs{}
			buildkite := services.NewBuildkite()
			run(RunParams{label, key, buildkiteJobId, autoRevertable, args}, osFs, buildkite)
		},
	}
)

func init() {
	runCmd.Flags().StringVarP(&label, "label", "", "", "the label of the step (required)")
	runCmd.Flags().StringVarP(&key, "key", "", "", "the key of the step (required)")
	runCmd.Flags().BoolVarP(&autoRevertable, "auto-revertable", "", false, "whether the step is auto-revertable (default=false)")
	runCmd.MarkFlagsRequiredTogether("label", "key")
}

func run(params RunParams, fs services.IFileService, buildkite services.IBuildkite) {
	if len(params.commands) == 0 {
		panic("need commands to run")
	}
	buildkite.LogSection(fmt.Sprintf("Running %s", params.label), false)
	startTime := time.Now()
	commands := params.commands
	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	exitCode := 0
	icon := ":bk-status-passed:"
	err := cmd.Run()
	if err != nil {
		icon = ":bk-status-failed:"
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	}
	elapsedTime := time.Since(startTime).Seconds()
	buildkite.LogSection(fmt.Sprintf("%s Finished %s in %.2f seconds\n", icon, params.label, elapsedTime), false)
	status := Status{label: params.label, key: params.key, exitCode: exitCode, autoRevertable: params.autoRevertable}
	statusManager := newStatusManager(params.jobId, fs)
	statusManager.append(status)
}
