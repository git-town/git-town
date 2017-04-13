package steps

import "github.com/Originate/git-town/lib/git"

// DeleteAncestorBranchesStep removes all ancestor information
// for the current branch
// from the Git Town configuration.
type DeleteAncestorBranchesStep struct {
	NoExpectedError
	NoUndoStep
}

// Run executes this step.
func (step DeleteAncestorBranchesStep) Run() error {
	git.DeleteAllAncestorBranches()
	return nil
}
