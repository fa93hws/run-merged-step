package cmd

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type StatusManagerTestSuite struct {
	suite.Suite
	simpleStatus    []Status
	realStatusBytes []byte

	mockedFs *MockedFs
}

func (suite *StatusManagerTestSuite) SetupSuite() {
	currentDir, _ := os.Getwd()
	realStatusFilePath := path.Join(currentDir, "fixtures", "simple.json")
	realStatusBytes, err := os.ReadFile(realStatusFilePath)
	if err != nil {
		panic(err)
	}
	suite.realStatusBytes = realStatusBytes
	err = json.Unmarshal(realStatusBytes, &suite.simpleStatus)
	if err != nil {
		panic(err)
	}
}

func (suite *StatusManagerTestSuite) SetupTest() {
	suite.mockedFs = &MockedFs{}
	suite.mockedFs.On("MkdirAll", mock.Anything, mock.Anything).Return(nil)
	suite.mockedFs.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.mockedFs.On("ReadFile", mock.Anything).Return(suite.realStatusBytes, nil)
	suite.mockedFs.On("RemoveAll", mock.Anything).Return(nil)
}

func (suite *StatusManagerTestSuite) TestStatusFilePath() {
	manager := NewStatusManager("/tmp", "job-id", suite.mockedFs)
	assert.Equal(suite.T(), "/tmp/job-id/merged_step_status.json", manager.filePath)
}

func (suite *StatusManagerTestSuite) TestMkDir() {
	manager := StatusManager{"/tmp/job-id/merged_step_status.json", suite.mockedFs}
	manager.mkdir()
	suite.mockedFs.AssertCalled(suite.T(), "MkdirAll", "/tmp/job-id", os.ModePerm)
}

func (suite *StatusManagerTestSuite) TestWriteEmptyStatus() {
	manager := StatusManager{"/tmp/job-id/merged_step_status.json", suite.mockedFs}
	emptyStatuses := []Status{}
	manager.writeToFile(emptyStatuses)
	suite.mockedFs.AssertCalled(suite.T(), "WriteFile", "/tmp/job-id/merged_step_status.json", []byte("[]"), os.ModePerm)
}

func (suite *StatusManagerTestSuite) TestWriteStatus() {
	manager := StatusManager{"/tmp/job-id/merged_step_status.json", suite.mockedFs}
	statuses := []Status{{
		Label:          "label-0",
		Key:            "key-0",
		ExitCode:       0,
		AutoRevertable: false,
	}}
	expectedJsonStr := `[
  {
    "label": "label-0",
    "key": "key-0",
    "exitCode": 0,
    "autoRevertable": false
  }
]`
	manager.writeToFile(statuses)
	suite.mockedFs.AssertCalled(suite.T(), "WriteFile", "/tmp/job-id/merged_step_status.json", mock.MatchedBy(func(jsonBytes []byte) bool {
		receivedJsonStr := string(jsonBytes)
		assert.Equal(suite.T(), expectedJsonStr, receivedJsonStr)
		return true
	}), os.ModePerm)
}

func (suite *StatusManagerTestSuite) TestReadStatus() {
	manager := StatusManager{"/tmp/job-id/merged_step_status.json", suite.mockedFs}
	statuses := manager.Read()
	assert.Equal(suite.T(), suite.simpleStatus, statuses)
	suite.mockedFs.AssertCalled(suite.T(), "ReadFile", "/tmp/job-id/merged_step_status.json")
}

func (suite *StatusManagerTestSuite) TestAppendStatus() {
	newStatus := Status{
		Label:          "l",
		Key:            "k",
		ExitCode:       0,
		AutoRevertable: false,
	}
	expectedStatus := append(suite.simpleStatus, newStatus)
	expectedJsonBytes, _ := json.MarshalIndent(expectedStatus, "", "  ")
	manager := StatusManager{"/tmp/job-id/merged_step_status.json", suite.mockedFs}
	manager.append(Status{
		Label:          "l",
		Key:            "k",
		ExitCode:       0,
		AutoRevertable: false,
	})
	suite.mockedFs.AssertCalled(suite.T(), "WriteFile", "/tmp/job-id/merged_step_status.json", expectedJsonBytes, os.ModePerm)
}

func (suite *StatusManagerTestSuite) TestRemove() {
	manager := StatusManager{"/tmp/job-id/merged_step_status.json", suite.mockedFs}
	manager.remove()
	suite.mockedFs.AssertCalled(suite.T(), "RemoveAll", "/tmp/job-id/merged_step_status.json")
}

func TestStatusManagerTestSuite(t *testing.T) {
	suite.Run(t, new(StatusManagerTestSuite))
}
