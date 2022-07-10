package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var buildkiteJobId string
var rootCmd = &cobra.Command{
	Use:   "run-merged-step",
	Short: "to execute merged step",
	Run: func(cmd *cobra.Command, args []string) {
		color.Red("need to provide command, run with --help to see a list of them")
		os.Exit(1)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&buildkiteJobId, "buildkite-job-id", "", "buildkite job id (required)")
	rootCmd.MarkPersistentFlagRequired("buildkite-job-id")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
