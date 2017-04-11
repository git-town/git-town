package steps

import (
	"github.com/Originate/git-town/lib/git"
)

type MergeTrackingBranchStep struct {
	NoAutomaticAbortOnError
}

func (step MergeTrackingBranchStep) CreateAbortStep() Step {
	return AbortMergeBranchStep{}
}

func (step MergeTrackingBranchStep) CreateContinueStep() Step {
	return ContinueMergeBranchStep{}
}

func (step MergeTrackingBranchStep) CreateUndoStep() Step {
	return ResetToShaStep{Hard: true, Sha: git.GetCurrentSha()}
}

func (step MergeTrackingBranchStep) Run() error {
	branchName := git.GetCurrentBranchName()
	if git.HasTrackingBranch(branchName) {
		return MergeBranchStep{BranchName: git.GetTrackingBranchName(branchName)}.Run()
	}
	return nil
}
