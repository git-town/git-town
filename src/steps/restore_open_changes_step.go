package steps

import (
	"errors"

	"github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/src/hosting"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	EmptyStep
}

func (step *RestoreOpenChangesStep) CreateUndoStep(backend *git.BackendCommands) (Step, error) {
	return &StashOpenChangesStep{}, nil
}

func (step *RestoreOpenChangesStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	err := run.Frontend.PopStash()
	if err != nil {
		return errors.New("conflicts between your uncommmitted changes and the main branch")
	}
	return nil
}
