package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PrepareTestSuit struct {
	suite.Suite
}

func (suite *PrepareTestSuit) TestPrepareToWriteEmptyFile() {
	mockFs := &MockedFs{}
	mockFs.On("TempDir").Return("/tmp")
	mockFs.On("MkdirAll", mock.Anything, mock.Anything).Return(nil)
	mockFs.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockBuildkite := &MockedBuildkite{}
	mockBuildkite.On("LogSection", mock.Anything, mock.Anything)
	prepare("job-id", mockFs, mockBuildkite)
	mockFs.AssertCalled(suite.T(), "TempDir")
	mockFs.AssertCalled(suite.T(), "MkdirAll", "/tmp/job-id", os.ModePerm)
	mockFs.AssertCalled(suite.T(), "WriteFile", "/tmp/job-id/merged_step_status.json", []byte("[]"), os.ModePerm)
}

func TestPrepareTestSuite(t *testing.T) {
	suite.Run(t, new(PrepareTestSuit))
}
