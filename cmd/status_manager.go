package cmd

import (
	"encoding/json"
	"os"
	"path"

	"github.com/fa93hws/run-merged-step/services"
)

var statusFileName = "merged_step_status.json"

type Status struct {
	Label          string `json:"label"`
	Key            string `json:"key"`
	ExitCode       int    `json:"exitCode"`
	AutoRevertable bool   `json:"autoRevertable"`
}

type IStatusManager interface {
	mkdir()
	writeToFile(status []Status)
	GetFilePath() string
	append(status Status)
}

type StatusManager struct {
	filePath string
	fs       services.IFileService
}

func NewStatusManager(buildkiteJobId string, fs services.IFileService) *StatusManager {
	filePath := path.Join(fs.TempDir(), buildkiteJobId, statusFileName)
	return &StatusManager{filePath, fs}
}

func (s *StatusManager) GetFilePath() string {
	return s.filePath
}

func (s *StatusManager) mkdir() {
	dir := path.Dir(s.filePath)
	err := s.fs.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (s *StatusManager) writeToFile(status []Status) {
	data, err := json.Marshal(status)
	if err != nil {
		panic(err)
	}
	err = s.fs.WriteFile(s.filePath, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (s *StatusManager) Read() []Status {
	data, err := s.fs.ReadFile(s.filePath)
	if err != nil {
		panic(err)
	}
	var statuses []Status
	json.Unmarshal(data, &statuses)
	return statuses
}

func (s *StatusManager) append(status Status) {
	statuses := s.Read()
	statuses = append(statuses, status)
	s.writeToFile(statuses)
}
