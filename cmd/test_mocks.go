package cmd

import (
	"os"

	"github.com/stretchr/testify/mock"
)

// status manager
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

func (m *MockedStatusManager) append(status Status) {
	m.Called(status)
}

func (m *MockedStatusManager) Read() []Status {
	args := m.Called()
	return args.Get(0).([]Status)
}

// fs
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

func (m *MockedFs) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	return args.Get(0).([]byte), args.Error(1)
}

// exec
type MockedExecService struct {
	mock.Mock
}

func (m *MockedExecService) Run(program string, programArgs []string, cwd *string, env *map[string]string) int {
	args := m.Called(program, programArgs, cwd, env)
	return args.Int(0)
}

// logger
type MockedLogger struct {
	mock.Mock
}

func (m *MockedLogger) LogSection(text string, collapse bool) {
	m.Called(text, collapse)
}
