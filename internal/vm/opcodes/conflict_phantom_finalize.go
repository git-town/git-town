package opcodes

import "github.com/git-town/git-town/v16/internal/vm/shared"

type ConflictPhantomFinalize struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictPhantomFinalize) Run(args shared.RunArgs) error {
	// See if we can commit now.
	// We can commit if all unresolved merge conflicts have been resolved.
	// Unresolved merge conflicts can remain if there are non-phantom merge conflicts.
	return nil
}
