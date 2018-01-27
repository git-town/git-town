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

	branchSha string
}

// AddUndoSteps adds the undo steps for this step to the undo step list
func (step *DeleteRemoteBranchStep) AddUndoSteps(stepList *StepList) {
	if step.IsTracking {
		stepList.Prepend(&CreateTrackingBranchStep{BranchName: step.BranchName})
	} else {
		stepList.Prepend(&CreateRemoteBranchStep{BranchName: step.BranchName, Sha: step.branchSha})
	}
}

// Run executes this step.
func (step *DeleteRemoteBranchStep) Run() error {
	if !step.IsTracking {
		step.branchSha = git.GetBranchSha(git.GetTrackingBranchName(step.BranchName))
	}
	return script.RunCommand("git", "push", "origin", ":"+step.BranchName)
}
