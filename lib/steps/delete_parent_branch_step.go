package steps

import "github.com/Originate/git-town/lib/config"

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
	parent := config.GetParentBranch(step.BranchName)
	if parent == "" {
		return NoOpStep{}
	} else {
		return SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: parent}
	}
}

func (step DeleteParentBranchStep) GetAutomaticAbortErrorMessage() string {
	return ""
}

func (step DeleteParentBranchStep) Run() error {
	config.DeleteParentBranch(step.BranchName)
	return nil
}

func (step DeleteParentBranchStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
