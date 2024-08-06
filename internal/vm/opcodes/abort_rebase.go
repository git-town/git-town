package opcodes

import "github.com/git-town/git-town/v14/internal/vm/shared"

// AbortRebase represents aborting on ongoing merge conflict.
// This opcode is used in the abort scripts for Git Town commands.
type AbortRebase struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *AbortRebase) Run(args shared.RunArgs) error {
	return args.Git.AbortRebase(args.Frontend)
}
