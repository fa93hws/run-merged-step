package services

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type ILogger interface {
	LogSection(string)
	LogCollapsedSection(string)
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

func (b Logger) LogSection(text string) {
	sprintf := color.New(color.Bold, color.FgMagenta).SprintfFunc()
	if !b.isOnBuildkite {
		b.LogInfo(sprintf(text))
		return
	}
	b.LogInfo(fmt.Sprintf("+++ %s", sprintf("%s", text)))
}

func (b Logger) LogCollapsedSection(text string) {
	sprintf := color.New(color.Bold, color.FgMagenta).SprintfFunc()
	if !b.isOnBuildkite {
		b.LogInfo(sprintf(text))
		return
	}
	b.LogInfo(fmt.Sprintf("--- %s", sprintf("%s", text)))
}

func (b Logger) LogInfo(text string) {
	fmt.Fprintln(os.Stderr, text)
}

type FakeLogger struct{}

func (f FakeLogger) LogSection(string)          {}
func (f FakeLogger) LogCollapsedSection(string) {}
func (f FakeLogger) LogInfo(string)             {}
