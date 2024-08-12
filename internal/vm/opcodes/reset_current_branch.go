package opcodes

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// ResetCurrentBranch resets all commits in the current branch.
type ResetCurrentBranch struct {
	Base                    gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ResetCurrentBranch) Run(args shared.RunArgs) error {
	return args.Git.ResetBranch(args.Frontend, self.Base)
}
