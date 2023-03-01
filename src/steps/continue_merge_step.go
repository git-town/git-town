package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// ContinueMergeStep finishes an ongoing merge conflict
// assuming all conflicts have been resolved by the user.
type ContinueMergeStep struct {
	EmptyStep
}

func (step *ContinueMergeStep) CreateContinueStep() Step {
	return step
}

func (step *ContinueMergeStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	hasMergeInprogress, err := repo.Silent.HasMergeInProgress()
	if err != nil {
		return err
	}
	if hasMergeInprogress {
		return repo.Logging.CommitNoEdit()
	}
	return nil
}
