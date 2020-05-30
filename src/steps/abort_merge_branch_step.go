package steps

import "github.com/git-town/git-town/src/git"

// AbortMergeBranchStep aborts the current merge conflict.
type AbortMergeBranchStep struct {
	NoOpStep
}

// Run executes this step.
func (step *AbortMergeBranchStep) Run(repo *git.ProdRepo) error {
	return repo.Logging.AbortMerge()
}
