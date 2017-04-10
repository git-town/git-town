package steps

import "github.com/Originate/git-town/lib/git"

type SetParentBranchStep struct {
	NoAutomaticAbort
	BranchName       string
	ParentBranchName string
}

func (step SetParentBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step SetParentBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step SetParentBranchStep) CreateUndoStep() Step {
	oldParent := git.GetParentBranch(step.BranchName)
	if oldParent == "" {
		return DeleteParentBranchStep{BranchName: step.BranchName}
	} else {
		return SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: oldParent}
	}
}

func (step SetParentBranchStep) Run() error {
	git.SetParentBranch(step.BranchName, step.ParentBranchName)
	return nil
}
