package opcodes

import "github.com/git-town/git-town/v15/internal/vm/shared"

// AbortMerge aborts the current merge conflict.
type AbortMerge struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *AbortMerge) Run(args shared.RunArgs) error {
	return args.Git.AbortMerge(args.Frontend)
}
