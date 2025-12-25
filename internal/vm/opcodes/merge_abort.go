package opcodes

import "github.com/git-town/git-town/v22/internal/vm/shared"

// MergeAbort aborts the current merge conflict.
type MergeAbort struct{}

func (self *MergeAbort) Run(args shared.RunArgs) error {
	return args.Git.AbortMerge(args.Frontend)
}
