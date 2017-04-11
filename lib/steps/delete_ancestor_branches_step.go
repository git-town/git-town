package steps

import "github.com/Originate/git-town/lib/git"

type DeleteAncestorBranchesStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
}

func (step DeleteAncestorBranchesStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step DeleteAncestorBranchesStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step DeleteAncestorBranchesStep) Run() error {
	git.DeleteAllAncestorBranches()
	return nil
}
