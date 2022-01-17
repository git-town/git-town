//nolint:ireturn
package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// ContinueRebaseBranchStep finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebaseBranchStep struct {
	NoOpStep
}

// CreateAbortStep returns the abort step for this step.
func (step *ContinueRebaseBranchStep) CreateAbortStep() Step {
	return &AbortRebaseBranchStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step *ContinueRebaseBranchStep) CreateContinueStep() Step {
	return step
}

// Run executes this step.
func (step *ContinueRebaseBranchStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	hasRebaseInProgress, err := repo.Silent.HasRebaseInProgress()
	if err != nil {
		return err
	}
	if hasRebaseInProgress {
		return repo.Logging.ContinueRebase()
	}
	return nil
}
