package steps

import "github.com/git-town/git-town/src/git"

// SetParentBranchStep registers the branch with the given name as a parent
// of the branch with the other given name.
type SetParentBranchStep struct {
	NoOpStep
	BranchName       string
	ParentBranchName string

	previousParent string
}

// CreateUndoStep returns the undo step for this step.
func (step *SetParentBranchStep) CreateUndoStep() Step {
	if step.previousParent == "" {
		return &DeleteParentBranchStep{BranchName: step.BranchName}
	}
	return &SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: step.previousParent}
}

// Run executes this step.
func (step *SetParentBranchStep) Run(repo *git.ProdRepo) error {
	step.previousParent = repo.GetParentBranch(step.BranchName)
	repo.SetParentBranch(step.BranchName, step.ParentBranchName)
	return nil
}
