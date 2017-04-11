package steps

import (
	"github.com/Originate/git-town/lib/script"
)

type PushTagsStep struct {
	NoAutomaticAbortOnError
}

func (step PushTagsStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step PushTagsStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step PushTagsStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step PushTagsStep) Run() error {
	return script.RunCommand("git", "push", "--tags")
}
