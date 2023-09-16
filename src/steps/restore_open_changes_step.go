package steps

import (
	"errors"

	"github.com/git-town/git-town/v9/src/messages"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	EmptyStep
}

func (step *RestoreOpenChangesStep) Run(args RunArgs) error {
	err := args.Runner.Frontend.PopStash()
	if err != nil {
		return errors.New(messages.DiffConflictWithMain)
	}
	return nil
}
