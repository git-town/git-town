package steps

import (
	"github.com/Originate/git-town/lib/script"
)

type DiscardOpenChangesStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
}

func (step DiscardOpenChangesStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step DiscardOpenChangesStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step DiscardOpenChangesStep) Run() error {
	return script.RunCommand("git", "reset", "--hard")
}
