package steps

import "github.com/Originate/git-town/lib/git"

// RemovePerennialBranch removes the branch with the given name as a perennial branch
type RemovePerennialBranch struct {
	NoAutomaticAbortOnError
	NoUndoStepAfterRun
	BranchName string
}

// CreateAbortStep returns the abort step for this step.
func (step RemovePerennialBranch) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step RemovePerennialBranch) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step RemovePerennialBranch) CreateUndoStepBeforeRun() Step {
	return AddPerennialBranch{BranchName: step.BranchName}
}

// Run executes this step.
func (step RemovePerennialBranch) Run() error {
	git.RemovePerennialBranch(step.BranchName)
	return nil
}
