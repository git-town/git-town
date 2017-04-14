package steps

import "github.com/Originate/git-town/lib/git"

// AddToPerennialBranches adds the branch with the given name as a perennial branch
type AddToPerennialBranches struct {
	NoOpStep
	BranchName string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step AddToPerennialBranches) CreateUndoStepBeforeRun() Step {
	return RemoveFromPerennialBranches{BranchName: step.BranchName}
}

// Run executes this step.
func (step AddToPerennialBranches) Run() error {
	git.AddToPerennialBranches(step.BranchName)
	return nil
}
