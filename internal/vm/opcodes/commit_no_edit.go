package opcodes

import "github.com/git-town/git-town/v16/internal/vm/shared"

type CommitNoEdit struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitNoEdit) Run(args shared.RunArgs) error {
	return args.Git.CommitNoEdit(args.Frontend)
}
