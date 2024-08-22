package opcodes

import (
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// RestoreOpenChanges restores stashed away changes into the workspace.
type DropStash struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *DropStash) Run(args shared.RunArgs) error {
	_ = args.Git.DropStash(args.Frontend)
	return nil
}
