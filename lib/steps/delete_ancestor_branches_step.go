package steps

import (
	"github.com/Originate/git-town/lib/config"
)

type DeleteAncestorBranchesStep struct{}

func (step DeleteAncestorBranchesStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step DeleteAncestorBranchesStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step DeleteAncestorBranchesStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step DeleteAncestorBranchesStep) Run() error {
	config.DeleteAllAncestorBranches()
	return nil
}
