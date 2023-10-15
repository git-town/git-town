package opcode

import "github.com/git-town/git-town/v9/src/vm/shared"

type StashOpenChanges struct {
	undeclaredOpcodeMethods
}

func (op *StashOpenChanges) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.Stash()
}
