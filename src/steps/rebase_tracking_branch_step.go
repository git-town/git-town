package steps

import (
	"github.com/Originate/git-town/src/git"
)

// RebaseTrackingBranchStep rebases the current branch against its tracking branch.
type RebaseTrackingBranchStep struct {
	NoOpStep
}

// CreateAbortStep returns the abort step for this step.
func (step RebaseTrackingBranchStep) CreateAbortStep() Step {
	return AbortRebaseBranchStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step RebaseTrackingBranchStep) CreateContinueStep() Step {
	return ContinueRebaseBranchStep{}
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step RebaseTrackingBranchStep) CreateUndoStepBeforeRun() Step {
	return ResetToShaStep{Hard: true, Sha: git.GetCurrentSha()}
}

// Run executes this step.
func (step RebaseTrackingBranchStep) Run() error {
	branchName := git.GetCurrentBranchName()
	if git.HasTrackingBranch(branchName) {
		return RebaseBranchStep{BranchName: git.GetTrackingBranchName(branchName)}.Run()
	}
	return nil
}
