package services

import (
	"os"

	"github.com/fatih/color"
)

type ILogger interface {
	LogSection(string, bool)
}

type Logger struct {
	isOnBuildkite bool
}

func NewLogger() *Logger {
	env := os.Getenv("BUILDKITE")
	isOnBuildkite := env == "true"
	return &Logger{isOnBuildkite}
}

func (b Logger) LogSection(text string, collapsed bool) {
	sprintf := color.New(color.Bold, color.FgMagenta).SprintfFunc()
	if !b.isOnBuildkite {
		print(sprintf("%s\n", text))
		return
	}
	var section string
	if collapsed {
		section = "---"
	} else {
		section = "+++"
	}
	print(section + " " + sprintf("%s\n", text))
}

type FakeLogger struct{}

func (f FakeLogger) LogSection(text string, collapsed bool) {}
