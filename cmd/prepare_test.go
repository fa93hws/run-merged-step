package cmd

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockedStatusManager struct {
	mock.Mock
}

func (m *MockedStatusManager) mkdir() {
	m.Called()
}

func (m *MockedStatusManager) writeToFile(status []Status) {
	m.Called(status)
}

func (m *MockedStatusManager) GetFilePath() string {
	m.Called()
	return ""
}

type MockedBuildkite struct {
	mock.Mock
}

func (m *MockedBuildkite) LogSection(text string, collapse bool) {
	m.Called(text, collapse)
}

type PrepareTestSuit struct {
	suite.Suite
	mockedStatusManager *MockedStatusManager
	mockedBuildkite     *MockedBuildkite

	fakeMkDir    *mock.Call
	fakeWrite    *mock.Call
	fakeFilePath *mock.Call
}

func (suite *PrepareTestSuit) SetupSuite() {
	suite.mockedStatusManager = &MockedStatusManager{}
	suite.fakeMkDir = suite.mockedStatusManager.On("mkdir")
	suite.fakeWrite = suite.mockedStatusManager.On("writeToFile", mock.Anything)
	suite.fakeFilePath = suite.mockedStatusManager.On("GetFilePath")

	suite.mockedBuildkite = &MockedBuildkite{}
	suite.mockedBuildkite.On("LogSection", mock.Anything, mock.Anything)
}

func (suite *PrepareTestSuit) BeforeTest() {
	suite.fakeMkDir.Unset()
	suite.fakeWrite.Unset()
	suite.fakeFilePath.Unset()
}

func (suite *PrepareTestSuit) TestPrepareToWriteEmptyFile() {
	prepare(suite.mockedStatusManager, suite.mockedBuildkite)
	suite.mockedStatusManager.AssertCalled(suite.T(), "mkdir")
	suite.mockedStatusManager.AssertCalled(suite.T(), "writeToFile", []Status{})
}

func TestPrepareTestSuite(t *testing.T) {
	suite.Run(t, new(PrepareTestSuit))
}
