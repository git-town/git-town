package opcode

import (
	"errors"

	"github.com/git-town/git-town/v11/src/messages"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// RestoreOpenChanges restores stashed away changes into the workspace.
type RestoreOpenChanges struct {
	undeclaredOpcodeMethods
}

func (self *RestoreOpenChanges) Run(args shared.RunArgs) error {
	stashSize, err := args.Runner.Backend.StashSnapshot()
	if err != nil {
		return err
	}
	if stashSize == 0 {
		return nil
	}
	err = args.Runner.Frontend.PopStash()
	if err != nil {
		return errors.New(messages.DiffConflictWithMain)
	}
	return nil
}
