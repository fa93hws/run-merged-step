package services

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type ILogger interface {
	LogSection(string, bool)
	LogInfo(string)
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
		fmt.Fprintln(os.Stderr, sprintf("%s\n", text))
		return
	}
	var section string
	if collapsed {
		section = "---"
	} else {
		section = "+++"
	}
	fmt.Fprintf(os.Stderr, "%s %s", section, sprintf("%s\n", text))
}

func (b Logger) LogInfo(text string) {
	fmt.Fprintln(os.Stderr, text)
}

type FakeLogger struct{}

func (f FakeLogger) LogSection(string, bool) {}
func (f FakeLogger) LogInfo(string)          {}
