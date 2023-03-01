package steps

import (
	"errors"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	EmptyStep
}

func (step *RestoreOpenChangesStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &StashOpenChangesStep{}, nil
}

func (step *RestoreOpenChangesStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	err := repo.Logging.PopStash()
	if err != nil {
		return errors.New("conflicts between your uncommmitted changes and the main branch")
	}
	return nil
}
