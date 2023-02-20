package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// ContinueMergeBranchStep finishes an ongoing merge conflict
// assuming all conflicts have been resolved by the user.
type ContinueMergeBranchStep struct {
	NoOpStep
}

func (step *ContinueMergeBranchStep) CreateContinueStep() Step {
	return step
}

func (step *ContinueMergeBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	hasMergeInprogress, err := repo.Silent.HasMergeInProgress()
	if err != nil {
		return err
	}
	if hasMergeInprogress {
		return repo.Logging.CommitNoEdit()
	}
	return nil
}
