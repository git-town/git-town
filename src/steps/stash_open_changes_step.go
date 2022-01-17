//nolint:ireturn
package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// StashOpenChangesStep stores all uncommitted changes on the Git stash.
type StashOpenChangesStep struct {
	NoOpStep
}

// CreateUndoStep returns the undo step for this step.
func (step *StashOpenChangesStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &RestoreOpenChangesStep{}, nil
}

// Run executes this step.
func (step *StashOpenChangesStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	return repo.Logging.Stash()
}
