package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

type StashOpenChangesStep struct {
	EmptyStep
}

func (step *StashOpenChangesStep) CreateUndoStep(backend *git.BackendCommands) (Step, error) {
	return &RestoreOpenChangesStep{}, nil
}

func (step *StashOpenChangesStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	return run.Frontend.Stash()
}
