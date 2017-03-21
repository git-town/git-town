package steps

import (
	"github.com/Originate/git-town/lib/config"
)

type DeleteAncestorBranches struct {
	BranchName string
}

func (step DeleteAncestorBranches) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step DeleteAncestorBranches) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step DeleteAncestorBranches) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step DeleteAncestorBranches) Run() error {
	config.DeleteAncestorBranches(step.BranchName)
	return nil
}
