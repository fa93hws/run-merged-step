package cmd

import (
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
	fs       services.FileService
}

func newStatusManager(buildkiteJobId string, fs services.FileService) *StatusManager {
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
