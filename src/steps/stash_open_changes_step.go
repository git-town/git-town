package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

type StashOpenChangesStep struct {
	EmptyStep
}

func (step *StashOpenChangesStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&RestoreOpenChangesStep{}}, nil
}

func (step *StashOpenChangesStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Frontend.Stash()
}
