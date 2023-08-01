package steps

import (
	"errors"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	EmptyStep
}

func (step *RestoreOpenChangesStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&StashOpenChangesStep{}}, nil
}

func (step *RestoreOpenChangesStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	err := run.Frontend.PopStash()
	if err != nil {
		return errors.New(messages.DiffConflictWithMain)
	}
	return nil
}
