package e2e

import (
	"os"
	"path"
	"testing"

	"github.com/fa93hws/run-merged-step/cmd"
	"github.com/fa93hws/run-merged-step/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	tempDir    string
	binPath    string
	passScript string
	failScript string
}

func (suite *E2ETestSuite) SetupSuite() {
	currentDir, _ := os.Getwd()
	suite.tempDir = path.Join(currentDir, "fixtures", "temp")
	suite.binPath = getBinaryPath()
	suite.passScript = path.Join(currentDir, "fixtures", "pass.sh")
	suite.failScript = path.Join(currentDir, "fixtures", "fail.sh")
}

func (suite *E2ETestSuite) TestAllCommandsPass() {
	jobId := "e2e-test-all-commands-pass"
	runCommand([]string{suite.binPath, "--buildkite-job-id", jobId, "--temp-dir", suite.tempDir, "prepare"})
	runCommand([]string{suite.binPath, "--buildkite-job-id", jobId, "--temp-dir", suite.tempDir, "run", "--label", "foo-label", "--key", "foo-key", "--auto-revertable", "--", suite.passScript})
	runCommand([]string{suite.binPath, "--buildkite-job-id", jobId, "--temp-dir", suite.tempDir, "run", "--label", "bar-label", "--key", "bar-key", "--", suite.passScript})

	fs := services.OsFs{}
	statusManager := cmd.NewStatusManager(os.TempDir(), jobId, fs)
	assert.True(suite.T(), exists(statusManager.GetFilePath()), "status file should exist")
	statuses := statusManager.Read()
	assert.Equal(suite.T(), statuses, []cmd.Status{{
		Label:          "foo-label",
		Key:            "foo-key",
		ExitCode:       0,
		AutoRevertable: true,
	}, {
		Label:          "bar-label",
		Key:            "bar-key",
		ExitCode:       0,
		AutoRevertable: false,
	},
	})
}

func TestPrepareTestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
