package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockedFs struct {
	mock.Mock
}

func (m *MockedFs) MkdirAll(name string, perm os.FileMode) error {
	args := m.Called(name, perm)
	return args.Error(0)
}

func (m *MockedFs) WriteFile(name string, data []byte, perm os.FileMode) error {
	args := m.Called(name, data, perm)
	return args.Error(0)
}

func (m *MockedFs) TempDir() string {
	args := m.Called()
	return args.String(0)
}

type MockedBuildkite struct {
	mock.Mock
}

func (m *MockedBuildkite) LogSection(text string, collapse bool) {
	m.Called(text, collapse)
}

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
