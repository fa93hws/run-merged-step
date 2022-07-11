package e2e

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	tempDir          string
	binPath          string
	passScript       string
	failScript       string
	autoRevertScript string
}

func (suite *E2ETestSuite) SetupSuite() {
	currentDir, _ := os.Getwd()
	fixtureDir := path.Join(currentDir, "fixtures")
	suite.tempDir = path.Join(fixtureDir, "temp")
	suite.binPath = getBinaryPath()
	suite.passScript = path.Join(fixtureDir, "pass_when_env_is_set.sh")
	suite.failScript = path.Join(fixtureDir, "fail.sh")
	suite.autoRevertScript = path.Join(fixtureDir, "upload_auto_revert_signal.sh")
	os.RemoveAll(suite.tempDir)
}

func (suite *E2ETestSuite) TestAllNonAutoRevertableCommandsPass() {
	jobId := "simple-pass"
	signalFilePath := path.Join(suite.tempDir, jobId, "signal.txt")
	commonArgs := []string{suite.binPath, "--buildkite-job-id", jobId, "--temp-dir", suite.tempDir}
	exitCode := runCommand(append(commonArgs, "prepare"), nil)
	assert.Equal(suite.T(), exitCode, 0, "prepare command should exit with 0")
	exitCode = runCommand(append(commonArgs, "run", "--label", "foo-label", "--key", "foo-key", "--", suite.passScript), nil)
	assert.Equal(suite.T(), exitCode, 0, "run command should always exit with 0")
	exitCode = runCommand(append(commonArgs, "run", "--label", "bar-label", "--key", "bar-key", "--", suite.passScript), nil)
	assert.Equal(suite.T(), exitCode, 0, "run command should always exit with 0")
	exitCode = runCommand(append(commonArgs, "report", "--auto-revertable-script", suite.autoRevertScript), &map[string]string{
		"AUTO_REVERT_OUTPUT_PATH": signalFilePath,
	})
	assert.Equal(suite.T(), exitCode, 0, "report should exit 0 when all commands pass")
	assert.False(suite.T(), exists(signalFilePath), "signal file should not be written when no command is auto-revertable")
}

func (suite *E2ETestSuite) TestSomeNonAutoRevertableCommandsFail() {
	jobId := "simple-fail"
	signalFilePath := path.Join(suite.tempDir, jobId, "signal.txt")
	commonArgs := []string{suite.binPath, "--buildkite-job-id", jobId, "--temp-dir", suite.tempDir}
	exitCode := runCommand(append(commonArgs, "prepare"), nil)
	assert.Equal(suite.T(), exitCode, 0, "prepare command should exit with 0")
	exitCode = runCommand(append(commonArgs, "run", "--label", "foo-label", "--key", "foo-key", "--", suite.failScript), nil)
	assert.Equal(suite.T(), exitCode, 0, "run command should always exit with 0")
	exitCode = runCommand(append(commonArgs, "run", "--label", "bar-label", "--key", "bar-key", "--", suite.passScript), nil)
	assert.Equal(suite.T(), exitCode, 0, "run command should always exit with 0")
	exitCode = runCommand(append(commonArgs, "report", "--auto-revertable-script", suite.autoRevertScript), &map[string]string{
		"AUTO_REVERT_OUTPUT_PATH": signalFilePath,
	})
	assert.Equal(suite.T(), exitCode, 1, "report should exit 1 when some command failed")
	assert.False(suite.T(), exists(signalFilePath), "signal file should not be written when no command is auto-revertable")
}

func (suite *E2ETestSuite) TestAllAutoRevertableCommandsPass() {
	jobId := "pass-auto-revertable"
	signalFilePath := path.Join(suite.tempDir, jobId, "signal.txt")
	commonArgs := []string{suite.binPath, "--buildkite-job-id", jobId, "--temp-dir", suite.tempDir}
	exitCode := runCommand(append(commonArgs, "prepare"), nil)
	assert.Equal(suite.T(), exitCode, 0, "prepare command should exit with 0")
	exitCode = runCommand(append(commonArgs, "run", "--label", "foo-label", "--key", "foo-key", "--auto-revertable", "--", suite.passScript), nil)
	assert.Equal(suite.T(), exitCode, 0, "run command should always exit with 0")
	exitCode = runCommand(append(commonArgs, "run", "--label", "bar-label", "--key", "bar-key", "--", suite.passScript), nil)
	assert.Equal(suite.T(), exitCode, 0, "run command should always exit with 0")
	exitCode = runCommand(append(commonArgs, "report", "--auto-revertable-script", suite.autoRevertScript), &map[string]string{
		"AUTO_REVERT_OUTPUT_PATH": signalFilePath,
	})
	assert.Equal(suite.T(), exitCode, 0, "report should exit 0 when all commands pass")
	assert.True(suite.T(), exists(signalFilePath), "signal file should be written if any command is auto-revertable")
	signal := readFileContent(signalFilePath)
	assert.Equal(suite.T(), signal, "passed")
}

func (suite *E2ETestSuite) TestSomeAutoRevertableCommandsFail() {
	jobId := "fail-auto-revertable"
	signalFilePath := path.Join(suite.tempDir, jobId, "signal.txt")
	commonArgs := []string{suite.binPath, "--buildkite-job-id", jobId, "--temp-dir", suite.tempDir}
	exitCode := runCommand(append(commonArgs, "prepare"), nil)
	assert.Equal(suite.T(), exitCode, 0, "prepare command should exit with 0")
	exitCode = runCommand(append(commonArgs, "run", "--label", "foo-label", "--key", "foo-key", "--auto-revertable", "--", suite.failScript), nil)
	assert.Equal(suite.T(), exitCode, 0, "run command should always exit with 0")
	exitCode = runCommand(append(commonArgs, "run", "--label", "bar-label", "--key", "bar-key", "--", suite.passScript), nil)
	assert.Equal(suite.T(), exitCode, 0, "run command should always exit with 0")
	exitCode = runCommand(append(commonArgs, "report", "--auto-revertable-script", suite.autoRevertScript), &map[string]string{
		"AUTO_REVERT_OUTPUT_PATH": signalFilePath,
	})
	assert.Equal(suite.T(), exitCode, 1, "report should exit 1 when some command failed")
	assert.True(suite.T(), exists(signalFilePath), "signal file should be written if any command is auto-revertable")
	signal := readFileContent(signalFilePath)
	assert.Equal(suite.T(), signal, "failed")
}

func (suite *E2ETestSuite) TestSomeNonAutoRevertableCommandsFailButAutoRevertableOnePass() {
	jobId := "fail-auto-revert-pass"
	signalFilePath := path.Join(suite.tempDir, jobId, "signal.txt")
	commonArgs := []string{suite.binPath, "--buildkite-job-id", jobId, "--temp-dir", suite.tempDir}
	exitCode := runCommand(append(commonArgs, "prepare"), nil)
	assert.Equal(suite.T(), exitCode, 0, "prepare command should exit with 0")
	exitCode = runCommand(append(commonArgs, "run", "--label", "foo-label", "--key", "foo-key", "--auto-revertable", "--", suite.passScript), nil)
	assert.Equal(suite.T(), exitCode, 0, "run command should always exit with 0")
	exitCode = runCommand(append(commonArgs, "run", "--label", "bar-label", "--key", "bar-key", "--", suite.failScript), nil)
	assert.Equal(suite.T(), exitCode, 0, "run command should always exit with 0")
	exitCode = runCommand(append(commonArgs, "report", "--auto-revertable-script", suite.autoRevertScript), &map[string]string{
		"AUTO_REVERT_OUTPUT_PATH": signalFilePath,
	})
	assert.Equal(suite.T(), exitCode, 1, "report should exit 1 when some command failed")
	assert.True(suite.T(), exists(signalFilePath), "signal file should be written if any command is auto-revertable")
	signal := readFileContent(signalFilePath)
	assert.Equal(suite.T(), signal, "passed")
}

func TestPrepareTestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
