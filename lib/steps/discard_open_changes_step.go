package steps

import (
	"github.com/Originate/git-town/lib/script"
)

type DiscardOpenChangesStep struct{}

func (step DiscardOpenChangesStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step DiscardOpenChangesStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step DiscardOpenChangesStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step DiscardOpenChangesStep) GetAutomaticAbortErrorMessage() string {
	return ""
}

func (step DiscardOpenChangesStep) Run() error {
	return script.RunCommand("git", "reset", "--hard")
}

func (step DiscardOpenChangesStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
