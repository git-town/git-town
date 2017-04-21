package steps

import "github.com/Originate/git-town/lib/git"

// RemoveFromPerennialBranches removes the branch with the given name as a perennial branch
type RemoveFromPerennialBranches struct {
	NoOpStep
	BranchName string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step RemoveFromPerennialBranches) CreateUndoStepBeforeRun() Step {
	return AddToPerennialBranches{BranchName: step.BranchName}
}

// Run executes this step.
func (step RemoveFromPerennialBranches) Run() error {
	git.RemoveFromPerennialBranches(step.BranchName)
	return nil
}
