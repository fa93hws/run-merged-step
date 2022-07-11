package cmd

import (
	"fmt"
	"time"

	"github.com/fa93hws/run-merged-step/services"
	"github.com/spf13/cobra"
)

type RunParams struct {
	label          string
	key            string
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
			logger := services.NewLogger()
			execService := services.NewExecService(nil)
			osFs := &services.OsFs{}
			statusManager := NewStatusManager(tempDir, buildkiteJobId, osFs)
			run(RunParams{label, key, autoRevertable, args}, statusManager, logger, execService)
		},
	}
)

func init() {
	runCmd.Flags().StringVarP(&label, "label", "", "", "the label of the step (required)")
	runCmd.Flags().StringVarP(&key, "key", "", "", "the key of the step (required)")
	runCmd.Flags().BoolVarP(&autoRevertable, "auto-revertable", "", false, "whether the step is auto-revertable (default=false)")
	runCmd.MarkFlagsRequiredTogether("label", "key")
}

func run(params RunParams, statusManager IStatusManager, logger services.ILogger, exec services.IExecService) {
	if len(params.commands) == 0 {
		panic("need commands to run")
	}
	logger.LogSection(fmt.Sprintf("Running %s", params.label), false)
	startTime := time.Now()
	commands := params.commands
	exitCode := exec.Run(commands[0], commands[1:])

	var icon string
	if exitCode == 0 {
		icon = ":bk-status-passed:"
	} else {
		icon = ":bk-status-failed:"
	}
	elapsedTime := time.Since(startTime).Seconds()
	logger.LogSection(fmt.Sprintf("%s Finished %s in %.2f seconds\n", icon, params.label, elapsedTime), false)
	status := Status{Label: params.label, Key: params.key, ExitCode: exitCode, AutoRevertable: params.autoRevertable}
	statusManager.append(status)
}
