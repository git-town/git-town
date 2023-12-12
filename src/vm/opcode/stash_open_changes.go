package opcode

import "github.com/git-town/git-town/v11/src/vm/shared"

type StashOpenChanges struct {
	undeclaredOpcodeMethods
}

func (self *StashOpenChanges) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.Stash()
}
