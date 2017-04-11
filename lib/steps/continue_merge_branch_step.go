package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

type ContinueMergeBranchStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
}

func (step ContinueMergeBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step ContinueMergeBranchStep) CreateContinueStep() Step {
	return step
}

func (step ContinueMergeBranchStep) Run() error {
	if git.IsMergeInProgress() {
		return script.RunCommand("git", "commit", "--no-edit")
	}
	return nil
}
