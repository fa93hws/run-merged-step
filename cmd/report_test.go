package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func createDummyStatus(name string, exitCode int, autoRevertable bool) Status {
	return Status{
		Label:          fmt.Sprintf("%s-label", name),
		Key:            fmt.Sprintf("%s-key", name),
		ExitCode:       exitCode,
		AutoRevertable: autoRevertable,
	}
}

type ReportTestSuite struct {
	suite.Suite

	mockedStatusManager *MockedStatusManager
	mockedLogger        *MockedLogger
	mockedExecService   *MockedExecService
}

func (suite *ReportTestSuite) SetupTest() {
	suite.mockedStatusManager = &MockedStatusManager{}
	suite.mockedStatusManager.On("remove").Return()

	suite.mockedLogger = &MockedLogger{}
	suite.mockedLogger.On("LogSection", mock.Anything, mock.Anything).Return()
	suite.mockedLogger.On("LogInfo", mock.Anything).Return()

	suite.mockedExecService = &MockedExecService{}
	suite.mockedExecService.On("Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0)
}

func (suite *ReportTestSuite) TestAllNonAutoRevertableCommandsPass() {
	statuses := []Status{
		createDummyStatus("command1", 0, false),
		createDummyStatus("command2", 0, false),
		createDummyStatus("command3", 0, false),
	}
	suite.mockedStatusManager.On("Read").Return(statuses).Once()
	exitCode := report(suite.mockedStatusManager, "upload_auto_revert_signal", suite.mockedLogger, suite.mockedExecService)
	suite.Equal(0, exitCode)
	suite.mockedLogger.AssertCalled(suite.T(), "LogSection", ":bk-status-passed: All step passed", false)
	suite.mockedExecService.AssertNotCalled(suite.T(), "Run")
}

func (suite *ReportTestSuite) TestAllAutoRevertableCommandsPass() {
	statuses := []Status{
		createDummyStatus("command1", 0, false),
		createDummyStatus("command2", 0, true),
		createDummyStatus("command3", 0, false),
	}
	suite.mockedStatusManager.On("Read").Return(statuses).Once()
	exitCode := report(suite.mockedStatusManager, "upload_auto_revert_signal", suite.mockedLogger, suite.mockedExecService)
	suite.Equal(0, exitCode)
	suite.mockedLogger.AssertCalled(suite.T(), "LogSection", ":bk-status-passed: All step passed", false)
	suite.mockedExecService.AssertCalled(suite.T(), "Run", "upload_auto_revert_signal", []string{"passed"}, (*string)(nil), (*map[string]string)(nil))
}

func (suite *ReportTestSuite) TestSomeNonAutoRevertableCommandsFail() {
	statuses := []Status{
		createDummyStatus("command1", 0, false),
		createDummyStatus("command2", 1, false),
		createDummyStatus("command3", 2, false),
	}
	suite.mockedStatusManager.On("Read").Return(statuses).Once()
	exitCode := report(suite.mockedStatusManager, "upload_auto_revert_signal", suite.mockedLogger, suite.mockedExecService)
	suite.Equal(1, exitCode)
	suite.mockedLogger.AssertCalled(suite.T(), "LogSection", ":bk-status-failed: Some steps failed", false)
	suite.mockedExecService.AssertNotCalled(suite.T(), "Run")
}

func (suite *ReportTestSuite) TestSomeAutoRevertableCommandsFail() {
	statuses := []Status{
		createDummyStatus("command1", 0, false),
		createDummyStatus("command2", 1, true),
		createDummyStatus("command3", 2, false),
	}
	suite.mockedStatusManager.On("Read").Return(statuses).Once()
	exitCode := report(suite.mockedStatusManager, "upload_auto_revert_signal", suite.mockedLogger, suite.mockedExecService)
	suite.Equal(1, exitCode)
	suite.mockedLogger.AssertCalled(suite.T(), "LogSection", ":bk-status-failed: Some steps failed", false)
	suite.mockedExecService.AssertCalled(suite.T(), "Run", "upload_auto_revert_signal", []string{"failed"}, (*string)(nil), (*map[string]string)(nil))
}

func (suite *ReportTestSuite) TestSomeNonAutoRevertableCommandsFailButAutoRevertableCommandsPass() {
	statuses := []Status{
		createDummyStatus("command1", 1, false),
		createDummyStatus("command2", 0, true),
		createDummyStatus("command3", 2, false),
	}
	suite.mockedStatusManager.On("Read").Return(statuses).Once()
	exitCode := report(suite.mockedStatusManager, "upload_auto_revert_signal", suite.mockedLogger, suite.mockedExecService)
	suite.Equal(1, exitCode)
	suite.mockedLogger.AssertCalled(suite.T(), "LogSection", ":bk-status-failed: Some steps failed", false)
	suite.mockedExecService.AssertCalled(suite.T(), "Run", "upload_auto_revert_signal", []string{"passed"}, (*string)(nil), (*map[string]string)(nil))
}

func (suite *ReportTestSuite) TestRemoveStatusFile() {
	statuses := []Status{
		createDummyStatus("command1", 0, false),
	}
	suite.mockedStatusManager.On("Read").Return(statuses).Once()
	report(suite.mockedStatusManager, "upload_auto_revert_signal", suite.mockedLogger, suite.mockedExecService)
	suite.mockedStatusManager.AssertCalled(suite.T(), "remove")
}

func TestReportTestSuite(t *testing.T) {
	suite.Run(t, new(ReportTestSuite))
}
