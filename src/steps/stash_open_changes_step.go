package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

type StashOpenChangesStep struct {
	EmptyStep
}

func (step *StashOpenChangesStep) CreateUndoStep(repo *git.PublicRepo) (Step, error) {
	return &RestoreOpenChangesStep{}, nil
}

func (step *StashOpenChangesStep) Run(repo *git.PublicRepo, connector hosting.Connector) error {
	return repo.Stash()
}
