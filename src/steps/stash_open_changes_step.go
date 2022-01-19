package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

type StashOpenChangesStep struct {
	NoOpStep
}

func (step *StashOpenChangesStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	return &RestoreOpenChangesStep{}, nil
}

func (step *StashOpenChangesStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	return repo.Logging.Stash()
}
