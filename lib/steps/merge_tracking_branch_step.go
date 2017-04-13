package steps

import (
	"github.com/Originate/git-town/lib/git"
)

// MergeTrackingBranchStep merges the tracking branch of the current branch
// into the current branch.
type MergeTrackingBranchStep struct {
	NoOpStep
}

// CreateAbortStep returns the abort step for this step.
func (step MergeTrackingBranchStep) CreateAbortStep() Step {
	return AbortMergeBranchStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step MergeTrackingBranchStep) CreateContinueStep() Step {
	return ContinueMergeBranchStep{}
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step MergeTrackingBranchStep) CreateUndoStepBeforeRun() Step {
	return ResetToShaStep{Hard: true, Sha: git.GetCurrentSha()}
}

// Run executes this step.
func (step MergeTrackingBranchStep) Run() error {
	branchName := git.GetCurrentBranchName()
	if git.HasTrackingBranch(branchName) {
		return MergeBranchStep{BranchName: git.GetTrackingBranchName(branchName)}.Run()
	}
	return nil
}
