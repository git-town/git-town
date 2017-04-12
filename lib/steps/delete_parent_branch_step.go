package steps

import "github.com/Originate/git-town/lib/git"

// DeleteParentBranchStep removes the parent branch entry in the Git Town configuration.
type DeleteParentBranchStep struct {
	BranchName string
}

// CreateAbortStep returns the abort step for this step.
func (step DeleteParentBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step DeleteParentBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step DeleteParentBranchStep) CreateUndoStep() Step {
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
