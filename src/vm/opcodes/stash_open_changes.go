package opcodes

import "github.com/git-town/git-town/v14/src/vm/shared"

type StashOpenChanges struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *StashOpenChanges) Run(args shared.RunArgs) error {
	return args.Git.Stash(args.Frontend)
}
