package steps

import "github.com/Originate/git-town/lib/git"

type DeleteAncestorBranchesStep struct{}

// CreateAbortStep returns the abort step for this step.
func (step DeleteAncestorBranchesStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step DeleteAncestorBranchesStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step DeleteAncestorBranchesStep) CreateUndoStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step DeleteAncestorBranchesStep) Run() error {
	git.DeleteAllAncestorBranches()
	return nil
}
