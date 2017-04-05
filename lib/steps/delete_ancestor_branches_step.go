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

func (step DeleteAncestorBranchesStep) GetAutomaticAbortErrorMessage() string {
	return ""
}

func (step DeleteAncestorBranchesStep) Run() error {
	config.DeleteAllAncestorBranches()
	return nil
}

func (step DeleteAncestorBranchesStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
