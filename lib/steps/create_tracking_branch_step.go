package steps

import (
	"github.com/Originate/git-town/lib/script"
)

type CreateTrackingBranchStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
	BranchName string
}

func (step CreateTrackingBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step CreateTrackingBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step CreateTrackingBranchStep) Run() error {
	return script.RunCommand("git", "push", "-u", "origin", step.BranchName)
}
