package steps

import "github.com/Originate/git-town/lib/gitconfig"

type DeleteParentBranchStep struct {
	BranchName string
}

func (step DeleteParentBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step DeleteParentBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step DeleteParentBranchStep) CreateUndoStep() Step {
	parent := gitconfig.GetParentBranch(step.BranchName)
	if parent == "" {
		return NoOpStep{}
	} else {
		return SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: parent}
	}
}

func (step DeleteParentBranchStep) Run() error {
	gitconfig.DeleteParentBranch(step.BranchName)
	return nil
}
