package steps

import (
	"errors"

	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	NoOpStep
}

// CreateUndoStep returns the undo step for this step.
func (step *RestoreOpenChangesStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &StashOpenChangesStep{}, nil
}

// Run executes this step.
func (step *RestoreOpenChangesStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	err := repo.Logging.PopStash()
	if err != nil {
		return errors.New("conflicts between your uncommmitted changes and the main branch")
	}
	return nil
}
