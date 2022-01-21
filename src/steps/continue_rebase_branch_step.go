package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// ContinueRebaseBranchStep finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebaseBranchStep struct {
	NoOpStep
}

func (step *ContinueRebaseBranchStep) CreateAbortStep() Step { //nolint:ireturn
	return &AbortRebaseBranchStep{}
}

func (step *ContinueRebaseBranchStep) CreateContinueStep() Step { //nolint:ireturn
	return step
}

func (step *ContinueRebaseBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	hasRebaseInProgress, err := repo.Silent.HasRebaseInProgress()
	if err != nil {
		return err
	}
	if hasRebaseInProgress {
		return repo.Logging.ContinueRebase()
	}
	return nil
}
