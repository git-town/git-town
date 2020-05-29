package steps

import (
	"github.com/git-town/git-town/src/git"
)

// StashOpenChangesStep stores all uncommitted changes on the Git stash.
type StashOpenChangesStep struct {
	NoOpStep
}

// CreateUndoStep returns the undo step for this step.
func (step *StashOpenChangesStep) CreateUndoStep() Step {
	return &RestoreOpenChangesStep{}
}

// Run executes this step.
func (step *StashOpenChangesStep) Run(repo *git.ProdRepo) error {
	return repo.Logging.Stash()
}
