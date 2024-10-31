package opcodes

import "github.com/git-town/git-town/v16/internal/vm/shared"

type ConflictPhantomFinalize struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictPhantomFinalize) Run(args shared.RunArgs) error {
	// TODO
	return nil
}
