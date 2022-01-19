//nolint:ireturn
package steps

import (
	"errors"

	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	NoOpStep
}

func (step *RestoreOpenChangesStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &StashOpenChangesStep{}, nil
}

func (step *RestoreOpenChangesStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	err := repo.Logging.PopStash()
	if err != nil {
		return errors.New("conflicts between your uncommmitted changes and the main branch")
	}
	return nil
}
