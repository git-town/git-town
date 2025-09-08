package opcodes

import "github.com/git-town/git-town/v21/internal/vm/shared"

type UndoLastCommit struct {
}

func (self *UndoLastCommit) Run(args shared.RunArgs) error {
	return args.Git.UndoLastCommit(args.Frontend)
}
