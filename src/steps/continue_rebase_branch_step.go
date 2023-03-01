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

func (step *ContinueRebaseBranchStep) CreateAbortStep() Step {
	return &AbortRebaseStep{}
}

func (step *ContinueRebaseBranchStep) CreateContinueStep() Step {
	return step
}

func (step *ContinueRebaseBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	hasRebaseInProgress, err := repo.Silent.HasRebaseInProgress()
	if err != nil {
		return err
	}
	if hasRebaseInProgress {
		return repo.Logging.ContinueRebase()
	}
	return nil
}
