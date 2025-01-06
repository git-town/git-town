package opcodes

import "github.com/git-town/git-town/v17/internal/vm/shared"

type UndoLastCommit struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *UndoLastCommit) Run(args shared.RunArgs) error {
	return args.Git.UndoLastCommit(args.Frontend)
}
