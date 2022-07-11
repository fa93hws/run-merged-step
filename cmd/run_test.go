package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RunCommandTestSuite struct {
	suite.Suite

	mockedExecService *MockedExecService
	fakeRun           *mock.Call

	mockedStatusManager *MockedStatusManager
	fakeAppend          *mock.Call

	mockedLogger   *MockedLogger
	fakeLogSection *mock.Call
}

func (suite *RunCommandTestSuite) SetupSuite() {
	suite.mockedExecService = &MockedExecService{}
	suite.mockedStatusManager = &MockedStatusManager{}
	suite.mockedLogger = &MockedLogger{}
}

func (suite *RunCommandTestSuite) SetupTest() {
	suite.fakeAppend = suite.mockedStatusManager.On("append", mock.Anything)
	suite.fakeLogSection = suite.mockedLogger.On("LogSection", mock.Anything, mock.Anything)
	if suite.fakeRun != nil {
		suite.fakeRun.Unset()
	}
}

func (suite *RunCommandTestSuite) TestRunCommandExitZero() {
	suite.fakeRun = suite.mockedExecService.On("Run", mock.Anything, mock.Anything).Return(0)
	run(RunParams{
		label:          "foo-label",
		key:            "foo-key",
		autoRevertable: false,
		commands:       []string{"echo", "-n", "foo"},
	}, suite.mockedStatusManager, suite.mockedLogger, suite.mockedExecService)
	suite.mockedExecService.AssertCalled(suite.T(), "Run", "echo", []string{"-n", "foo"})
	suite.mockedLogger.AssertCalled(suite.T(), "LogSection", mock.MatchedBy(func(message string) bool {
		return strings.Contains(message, ":bk-status-passed:")
	}), false)
	suite.mockedStatusManager.AssertCalled(suite.T(), "append", Status{
		Label:          "foo-label",
		Key:            "foo-key",
		AutoRevertable: false,
		ExitCode:       0,
	})
}

func (suite *RunCommandTestSuite) TestRunCommandExitNonZero() {
	suite.fakeRun = suite.mockedExecService.On("Run", mock.Anything, mock.Anything).Return(3)
	run(RunParams{
		label:          "bar-label",
		key:            "bar-key",
		autoRevertable: true,
		commands:       []string{"echo", "-n", "foo"},
	}, suite.mockedStatusManager, suite.mockedLogger, suite.mockedExecService)
	suite.mockedExecService.AssertCalled(suite.T(), "Run", "echo", []string{"-n", "foo"})
	suite.mockedLogger.AssertCalled(suite.T(), "LogSection", mock.MatchedBy(func(message string) bool {
		return strings.Contains(message, ":bk-status-failed:")
	}), false)
	suite.mockedStatusManager.AssertCalled(suite.T(), "append", Status{
		Label:          "bar-label",
		Key:            "bar-key",
		AutoRevertable: true,
		ExitCode:       3,
	})
}

func TestRunCommandTestSuite(t *testing.T) {
	suite.Run(t, new(RunCommandTestSuite))
}
