package steps

import "github.com/Originate/git-town/lib/git"

// AddPerennialBranch adds the branch with the given name as a perennial branch
type AddPerennialBranch struct {
	NoAutomaticAbortOnError
	NoUndoStepAfterRun
	BranchName string
}

// CreateAbortStep returns the abort step for this step.
func (step AddPerennialBranch) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step AddPerennialBranch) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step AddPerennialBranch) CreateUndoStepBeforeRun() Step {
	return RemovePerennialBranch{BranchName: step.BranchName}
}

// Run executes this step.
func (step AddPerennialBranch) Run() error {
	git.AddPerennialBranch(step.BranchName)
	return nil
}
