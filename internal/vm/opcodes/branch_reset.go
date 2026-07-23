package opcodes

import (
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/vm/shared"
)

type BranchReset struct {
	Target gitdomain.BranchName
}

func (self *BranchReset) Run(args shared.RunArgs) error {
	return args.Git.ResetBranch(args.Frontend, self.Target)
}
