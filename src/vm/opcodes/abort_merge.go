package opcodes

import "github.com/git-town/git-town/v14/src/vm/shared"

// AbortMerge aborts the current merge conflict.
type AbortMerge struct {
	undeclaredOpcodeMethods
}

func (self *AbortMerge) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.AbortMerge()
}
