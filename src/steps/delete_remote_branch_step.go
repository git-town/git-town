package steps

import (
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
)

// DeleteRemoteBranchStep deletes the current branch from the origin remote.
type DeleteRemoteBranchStep struct {
	NoOpStep
	BranchName string
	IsTracking bool

	branchSha string
}

// CreateUndoStep returns the undo step for this step.
func (step *DeleteRemoteBranchStep) CreateUndoStep() Step {
	if step.IsTracking {
		return &CreateTrackingBranchStep{BranchName: step.BranchName}
	}
	return &CreateRemoteBranchStep{BranchName: step.BranchName, Sha: step.branchSha}
}

// Run executes this step.
func (step *DeleteRemoteBranchStep) Run() error {
	if !step.IsTracking {
		step.branchSha = git.GetBranchSha(git.GetTrackingBranchName(step.BranchName))
	}
	return script.RunCommand("git", "push", "origin", ":"+step.BranchName)
}
