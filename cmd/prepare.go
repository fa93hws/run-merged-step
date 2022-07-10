package cmd

import (
	"github.com/fa93hws/run-merged-step/services"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(prepareCmd)
}

var prepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Preparation work for running merged step",
	Run: func(cmd *cobra.Command, args []string) {
		osFs := &services.OsFs{}
		prepare(buildkiteJobId, osFs)
	},
}

func prepare(jobId string, fs services.FileService) {
	statusManager := newStatusManager(jobId, fs)
	statusManager.mkdir()
	statusManager.writeToFile([]Status{})
}
