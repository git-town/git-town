package opcode

import (
	"errors"

	"github.com/git-town/git-town/v9/src/messages"
)

// RestoreOpenChanges restores stashed away changes into the workspace.
type RestoreOpenChanges struct {
	undeclaredOpcodeMethods
}

func (step *RestoreOpenChanges) Run(args RunArgs) error {
	err := args.Runner.Frontend.PopStash()
	if err != nil {
		return errors.New(messages.DiffConflictWithMain)
	}
	return nil
}
