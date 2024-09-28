package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// ResetCurrentBranch resets all commits in the current branch.
type ResetBranch struct {
	Target                  gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ResetBranch) Run(args shared.RunArgs) error {
	return args.Git.ResetBranch(args.Frontend, self.Target)
}
