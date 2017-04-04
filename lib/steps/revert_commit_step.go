package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

type RevertCommitStep struct{}

func (step RevertCommitStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step RevertCommitStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step RevertCommitStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step RevertCommitStep) Run() error {
	err := script.RunCommand("git", "revert", "HEAD")
	if err != nil {
		return err
	}
	return PushBranchStep{BranchName: git.GetCurrentBranchName()}.Run()
}
