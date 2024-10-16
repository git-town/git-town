package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// BranchCurrentReset resets all commits in the current branch.
type BranchCurrentReset struct {
	Base                    gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchCurrentReset) Run(args shared.RunArgs) error {
	return args.Git.ResetBranch(args.Frontend, self.Base)
}
