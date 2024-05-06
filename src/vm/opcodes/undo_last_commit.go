package opcodes

import "github.com/git-town/git-town/v14/src/vm/shared"

type UndoLastCommit struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *UndoLastCommit) Run(args shared.RunArgs) error {
	return args.Frontend.UndoLastCommit()
}
