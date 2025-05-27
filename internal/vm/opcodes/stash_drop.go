package opcodes

import (
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// RestoreOpenChanges restores stashed away changes into the workspace.
type StashDrop struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *StashDrop) Run(args shared.RunArgs) error {
	_ = args.Git.DropMostRecentStash(args.Frontend)
	return nil
}
