package steps

import "github.com/git-town/git-town/src/git"

// AddToPerennialBranches adds the branch with the given name as a perennial branch
type AddToPerennialBranches struct {
	NoOpStep
	BranchName string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step *AddToPerennialBranches) CreateUndoStepBeforeRun() Step {
	return &RemoveFromPerennialBranches{BranchName: step.BranchName}
}

// Run executes this step.
func (step *AddToPerennialBranches) Run() error {
	git.Config().AddToPerennialBranches(step.BranchName)
	return nil
}
