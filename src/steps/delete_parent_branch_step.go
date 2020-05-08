package steps

import "github.com/git-town/git-town/src/git"

// DeleteParentBranchStep removes the parent branch entry in the Git Town configuration.
type DeleteParentBranchStep struct {
	NoOpStep
	BranchName string

	previousParent string
}

// CreateUndoStep returns the undo step for this step.
func (step *DeleteParentBranchStep) CreateUndoStep() Step {
	if step.previousParent == "" {
		return &NoOpStep{}
	}
	return &SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: step.previousParent}
}

// Run executes this step.
func (step *DeleteParentBranchStep) Run() error {
	step.previousParent = git.Config().GetParentBranch(step.BranchName)
	git.Config().DeleteParentBranch(step.BranchName)
	return nil
}
