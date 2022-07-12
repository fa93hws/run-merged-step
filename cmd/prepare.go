package cmd

import (
	"github.com/fa93hws/run-merged-step/services"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var prepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Preparation work for running merged step",
	Run: func(cmd *cobra.Command, args []string) {
		statusManager := NewStatusManager(tempDir, buildkiteJobId, &services.OsFs{})
		logger := services.NewLogger()
		prepare(statusManager, logger)
	},
}

func prepare(statusManager IStatusManager, logger services.ILogger) {
	logger.LogSection("Creating status file")
	statusManager.mkdir()
	statusManager.writeToFile([]Status{})
	color.Green("Status file created at %s", statusManager.GetFilePath())
}
