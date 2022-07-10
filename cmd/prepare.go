package cmd

import (
	"os"
	"path"

	"github.com/fa93hws/run-merged-step/services"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(prepareCmd)
}

var prepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Preparation work for running merged step",
	Run: func(cmd *cobra.Command, args []string) {
		osFs := services.OsFs{}
		prepare(buildkiteJobId, osFs)
	},
}

func prepare(buildkiteJobId string, fs services.FileSystem) {
	tmpDir := path.Join(fs.TempDir(), buildkiteJobId)
	err := fs.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	statusFile := path.Join(tmpDir, statusFileName)
	err = fs.WriteFile(statusFile, []byte("[]"), os.ModePerm)
	if err != nil {
		panic(err)
	}
	color.Green("status file written to %v\n", statusFile)
}
