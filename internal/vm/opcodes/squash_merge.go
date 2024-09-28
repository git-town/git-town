package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type SquashMerge struct {
	Branch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *SquashMergeWorkflow) Run(args shared.RunArgs) error {
	return args.Git.SquashMerge(args.Frontend, self.Branch)
}
