package steps

import "github.com/Originate/git-town/lib/git"

type DeleteParentBranchStep struct {
	NoAutomaticAbortOnError
	NoUndoStepAfterRun
	BranchName string
}

func (step DeleteParentBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step DeleteParentBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step DeleteParentBranchStep) CreateUndoStepBeforeRun() Step {
	parent := git.GetParentBranch(step.BranchName)
	if parent == "" {
		return NoOpStep{}
	} else {
		return SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: parent}
	}
}

func (step DeleteParentBranchStep) Run() error {
	git.DeleteParentBranch(step.BranchName)
	return nil
}
