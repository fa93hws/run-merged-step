package cmd

import (
	"testing"

	"github.com/fa93hws/run-merged-step/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PrepareTestSuit struct {
	suite.Suite

	mockedStatusManager *MockedStatusManager
	fakeLogger          services.ILogger
}

func (suite *PrepareTestSuit) SetupSuite() {
	suite.mockedStatusManager = &MockedStatusManager{}
	suite.fakeLogger = &services.FakeLogger{}
}

func (suite *PrepareTestSuit) SetupTest() {
	suite.mockedStatusManager.On("mkdir")
	suite.mockedStatusManager.On("writeToFile", mock.Anything)
	suite.mockedStatusManager.On("GetFilePath")
}

func (suite *PrepareTestSuit) TestPrepareToWriteEmptyFile() {
	prepare(suite.mockedStatusManager, suite.fakeLogger)
	suite.mockedStatusManager.AssertCalled(suite.T(), "mkdir")
	suite.mockedStatusManager.AssertCalled(suite.T(), "writeToFile", []Status{})
}

func TestPrepareTestSuite(t *testing.T) {
	suite.Run(t, new(PrepareTestSuit))
}
