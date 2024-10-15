package opcodes

import "github.com/git-town/git-town/v16/internal/vm/shared"

// RebaseAbort represents aborting on ongoing merge conflict.
// This opcode is used in the abort scripts for Git Town commands.
type RebaseAbort struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseAbort) Run(args shared.RunArgs) error {
	return args.Git.AbortRebase(args.Frontend)
}
