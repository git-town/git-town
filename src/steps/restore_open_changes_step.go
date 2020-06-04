package steps

import (
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	NoOpStep
}

// CreateUndoStep returns the undo step for this step.
func (step *RestoreOpenChangesStep) CreateUndoStep() Step {
	return &StashOpenChangesStep{}
}

// Run executes this step.
func (step *RestoreOpenChangesStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	return repo.Logging.PopStash()
}
