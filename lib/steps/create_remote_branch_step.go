package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// CreateRemoteBranchStep pushes the current branch up to origin.
type CreateRemoteBranchStep struct {
	BranchName string
	Sha        string
}

// CreateAbortStep returns the abort step for this step.
func (step CreateRemoteBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step CreateRemoteBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step CreateRemoteBranchStep) CreateUndoStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step CreateRemoteBranchStep) Run() error {
	return script.RunCommand("git", "push", "origin", step.Sha+":refs/heads/"+step.BranchName)
}
