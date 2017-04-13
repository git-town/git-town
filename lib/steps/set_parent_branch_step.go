package steps

import "github.com/Originate/git-town/lib/git"

// SetParentBranchStep registers the branch with the given name as a parent
// of the branch with the other given name.
type SetParentBranchStep struct {
	NoExpectedError
	NoUndoStepAfterRun
	BranchName       string
	ParentBranchName string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step SetParentBranchStep) CreateUndoStepBeforeRun() Step {
	oldParent := git.GetParentBranch(step.BranchName)
	if oldParent == "" {
		return DeleteParentBranchStep{BranchName: step.BranchName}
	}
	return SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: oldParent}
}

// Run executes this step.
func (step SetParentBranchStep) Run() error {
	git.SetParentBranch(step.BranchName, step.ParentBranchName)
	return nil
}
