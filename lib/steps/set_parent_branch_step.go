package steps

import "github.com/Originate/git-town/lib/git"

// SetParentBranchStep registers the branch with the given name as a parent
// of the branch with the other given name.
type SetParentBranchStep struct {
	BranchName       string
	ParentBranchName string
}

// CreateAbortStep returns the abort step for this step.
func (step SetParentBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step SetParentBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step SetParentBranchStep) CreateUndoStep() Step {
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
