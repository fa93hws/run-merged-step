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
		osFs := &services.OsFs{}
		buildkite := services.NewBuildkite()
		prepare(buildkiteJobId, osFs, buildkite)
	},
}

func prepare(jobId string, fs services.IFileService, buildkite services.IBuildkite) {
	buildkite.LogSection("Creating status file", false)
	statusManager := newStatusManager(jobId, fs)
	statusManager.mkdir()
	statusManager.writeToFile([]Status{})
	color.Green("Status file created at %s", statusManager.filePath)
}
