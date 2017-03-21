package steps

import (
	"github.com/Originate/git-town/lib/config"
)

type SetParentBranchStep struct {
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
	return NoOpStep{}
}

func (step SetParentBranchStep) Run() error {
	config.SetParentBranch(step.BranchName, step.ParentBranchName)
	return nil
}
