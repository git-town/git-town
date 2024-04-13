package opcodes

import "github.com/git-town/git-town/v14/src/vm/shared"

type StashOpenChanges struct {
	undeclaredOpcodeMethods
}

func (self *StashOpenChanges) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.Stash()
}
