package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

// DeleteRemoteBranchStep deletes the current branch from the origin remote.
type DeleteRemoteBranchStep struct {
	BranchName string
	IsTracking bool
}

// CreateAbortStep returns the abort step for this step.
func (step DeleteRemoteBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step DeleteRemoteBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step DeleteRemoteBranchStep) CreateUndoStep() Step {
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
