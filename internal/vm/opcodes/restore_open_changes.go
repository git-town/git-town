package opcodes

import (
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// RestoreOpenChanges restores stashed away changes into the workspace.
type RestoreOpenChanges struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RestoreOpenChanges) Run(args shared.RunArgs) error {
	stashSize, err := args.Git.StashSize(args.Backend)
	if err != nil {
		return err
	}
	if stashSize == 0 && !args.Config.DryRun {
		return nil
	}
	args.PrependOpcodes(&PopStash{})
	return nil
}
