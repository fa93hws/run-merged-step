package cmd

import (
	"encoding/json"
	"os"
	"path"

	"github.com/fa93hws/run-merged-step/services"
)

var statusFileName = "merged_step_status.json"

type Status struct {
	label          string
	key            string
	exitCode       int
	autoRevertable bool
}

type StatusManager struct {
	filePath string
	fs       services.IFileService
}

func newStatusManager(buildkiteJobId string, fs services.IFileService) *StatusManager {
	filePath := path.Join(fs.TempDir(), buildkiteJobId, statusFileName)
	return &StatusManager{filePath, fs}
}

func (s *StatusManager) mkdir() {
	dir := path.Dir(s.filePath)
	err := s.fs.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (s *StatusManager) writeToFile(status []Status) {
	err := s.fs.WriteFile(s.filePath, []byte("[]"), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (s *StatusManager) read() []Status {
	data, err := s.fs.ReadFile(s.filePath)
	if err != nil {
		panic(err)
	}
	var statuses []Status
	json.Unmarshal(data, &statuses)
	return statuses
}

func (s *StatusManager) append(status Status) {
	statuses := s.read()
	statuses = append(statuses, status)
	s.writeToFile(statuses)
}