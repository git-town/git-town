package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// RestoreOpenChanges restores stashed away changes into the workspace.
type RestoreOpenChanges struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RestoreOpenChanges) Run(args shared.RunArgs) error {
	stashSize, err := args.Backend.StashSize()
	if err != nil {
		return err
	}
	if stashSize == 0 && !args.Config.DryRun {
		return nil
	}
	err = args.Frontend.PopStash()
	if err != nil {
		return errors.New(messages.DiffConflictWithMain)
	}
	return nil
}
