package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// RestoreOpenChanges restores stashed away changes into the workspace.
type RestoreOpenChanges struct {
	undeclaredOpcodeMethods
}

func (self *RestoreOpenChanges) Run(args shared.RunArgs) error {
	stashSize, err := args.Runner.Backend.StashSize()
	if err != nil {
		return err
	}
	if stashSize == 0 && !args.Runner.Config.DryRun {
		return nil
	}
	err = args.Runner.Frontend.PopStash()
	if err != nil {
		return errors.New(messages.DiffConflictWithMain)
	}
	return nil
}
