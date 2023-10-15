package opcode

import "github.com/git-town/git-town/v9/src/vm/shared"

// AbortRebase represents aborting on ongoing merge conflict.
// This step is used in the abort scripts for Git Town commands.
type AbortRebase struct {
	undeclaredOpcodeMethods
}

func (step *AbortRebase) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.AbortRebase()
}
