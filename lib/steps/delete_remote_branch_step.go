package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

type DeleteRemoteBranchStep struct {
	NoAutomaticAbortOnError
	BranchName string
	IsTracking bool
}

func (step DeleteRemoteBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step DeleteRemoteBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step DeleteRemoteBranchStep) CreateUndoStep() Step {
	if step.IsTracking {
		return CreateTrackingBranchStep{BranchName: step.BranchName}
	} else {
		sha := git.GetBranchSha(git.GetTrackingBranchName(step.BranchName))
		return CreateRemoteBranchStep{BranchName: step.BranchName, Sha: sha}
	}
}

func (step DeleteRemoteBranchStep) Run() error {
	return script.RunCommand("git", "push", "origin", ":"+step.BranchName)
}
