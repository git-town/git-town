package steps

import "github.com/git-town/git-town/src/git"

// AddToPerennialBranches adds the branch with the given name as a perennial branch
type AddToPerennialBranches struct {
	NoOpStep
	BranchName string
}

// CreateUndoStep returns the undo step for this step.
func (step *AddToPerennialBranches) CreateUndoStep() Step {
	return &RemoveFromPerennialBranches{BranchName: step.BranchName}
}

// Run executes this step.
func (step *AddToPerennialBranches) Run() error {
	git.Config().AddToPerennialBranches(step.BranchName)
	return nil
}
