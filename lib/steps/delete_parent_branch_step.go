package steps

import "github.com/Originate/git-town/lib/git"

// DeleteParentBranchStep removes the parent branch entry in the Git Town configuration.
type DeleteParentBranchStep struct {
	NoExpectedError
	NoUndoStepAfterRun
	BranchName string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step DeleteParentBranchStep) CreateUndoStepBeforeRun() Step {
	parent := git.GetParentBranch(step.BranchName)
	if parent == "" {
		return NoOpStep{}
	}
	return SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: parent}
}

// Run executes this step.
func (step DeleteParentBranchStep) Run() error {
	git.DeleteParentBranch(step.BranchName)
	return nil
}
