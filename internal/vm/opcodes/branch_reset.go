package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

type BranchReset struct {
	Target                  gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchReset) Run(args shared.RunArgs) error {
	return args.Git.ResetBranch(args.Frontend, self.Target)
}
