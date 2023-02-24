package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

type StashOpenChangesStep struct {
	NoOpStep
}

func (step *StashOpenChangesStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &RestoreOpenChangesStep{}, nil
}

func (step *StashOpenChangesStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return repo.Logging.Stash()
}
