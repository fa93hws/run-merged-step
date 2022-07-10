package services

import (
	"os"

	"github.com/fatih/color"
)

type IBuildkite interface {
	LogSection(string, bool)
}

type Buildkite struct {
	isOnBuildkite bool
}

func NewBuildkite() *Buildkite {
	env := os.Getenv("BUILDKITE")
	isOnBuildkite := env == "true"
	return &Buildkite{isOnBuildkite}
}

func (b Buildkite) LogSection(text string, collapsed bool) {
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
