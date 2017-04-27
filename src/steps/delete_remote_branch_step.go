package steps

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
)

// DeleteRemoteBranchStep deletes the current branch from the origin remote.
type DeleteRemoteBranchStep struct {
	NoOpStep
	BranchName string
	IsTracking bool
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step DeleteRemoteBranchStep) CreateUndoStepBeforeRun() Step {
	if step.IsTracking {
		return CreateTrackingBranchStep{BranchName: step.BranchName}
	}
	sha := git.GetBranchSha(git.GetTrackingBranchName(step.BranchName))
	return CreateRemoteBranchStep{BranchName: step.BranchName, Sha: sha}
}

// Run executes this step.
func (step DeleteRemoteBranchStep) Run() error {
	return script.RunCommand("git", "push", "origin", ":"+step.BranchName)
}
