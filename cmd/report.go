package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fa93hws/run-merged-step/services"
	"github.com/spf13/cobra"
)

var autoRevertableScriptPath string
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report the final status of a merged step",
	Run: func(cmd *cobra.Command, args []string) {
		statusManager := NewStatusManager(tempDir, buildkiteJobId, &services.OsFs{})
		execService := &services.ExecService{}
		logger := services.NewLogger()
		exitCode := report(statusManager, autoRevertableScriptPath, logger, execService)
		os.Exit(exitCode)
	},
}

func init() {
	reportCmd.Flags().StringVarP(&autoRevertableScriptPath, "auto-revertable-script", "", "", "path to upload auto revertable signal script (required)")
	runCmd.MarkFlagRequired("auto-revertable-script")
}

func filterFailedCommands(statuses []Status) []Status {
	result := []Status{}
	for _, status := range statuses {
		if status.ExitCode != 0 {
			result = append(result, status)
		}
	}
	return result
}

func hasAutoRevertable(statuses []Status) bool {
	for _, status := range statuses {
		if status.AutoRevertable {
			return true
		}
	}
	return false
}

func printResult(failedStatuses []Status, logger services.ILogger) {
	if len(failedStatuses) == 0 {
		logger.LogSection(":bk-status-passed: All step passed", false)
	}
	logger.LogSection(":bk-status-failed: Some steps failed", false)
	for _, failedStatus := range failedStatuses {
		logger.LogInfo(fmt.Sprintf("%s failed with exit code=%d\n", failedStatus.Label, failedStatus.ExitCode))
	}
}

func report(statusManager IStatusManager, autoRevertableScriptPath string, logger services.ILogger, exec services.IExecService) int {
	statuses := statusManager.Read()
	failed := filterFailedCommands(statuses)
	printResult(failed, logger)
	if hasAutoRevertable(statuses) {
		var signal string
		if hasAutoRevertable(failed) {
			signal = "failed"
		} else {
			signal = "passed"
		}
		exitCode := exec.Run(autoRevertableScriptPath, []string{signal}, nil, nil)
		if exitCode != 0 {
			panic(fmt.Errorf("failed to run auto revertable script %s, signal=%s", autoRevertableScriptPath, signal))
		}
	}

	logger.LogSection("Print status file", true)
	jsonBytes, _ := json.MarshalIndent(statuses, "", "  ")
	logger.LogInfo(string(jsonBytes))

	if len(failed) > 0 {
		return 1
	} else {
		return 0
	}
}
