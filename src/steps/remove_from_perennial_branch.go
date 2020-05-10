package steps

import "github.com/git-town/git-town/src/git"

// RemoveFromPerennialBranches removes the branch with the given name as a perennial branch
type RemoveFromPerennialBranches struct {
	NoOpStep
	BranchName string
}

// CreateUndoStep returns the undo step for this step.
func (step *RemoveFromPerennialBranches) CreateUndoStep() Step {
	return &AddToPerennialBranches{BranchName: step.BranchName}
}

// Run executes this step.
func (step *RemoveFromPerennialBranches) Run() error {
	git.Config().RemoveFromPerennialBranches(step.BranchName)
	return nil
}
