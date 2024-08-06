package opcodes

import "github.com/git-town/git-town/v15/internal/vm/shared"

type UndoLastCommit struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *UndoLastCommit) Run(args shared.RunArgs) error {
	return args.Git.UndoLastCommit(args.Frontend)
}
