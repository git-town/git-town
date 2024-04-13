package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// ResetCommitsInCurrentBranch resets all commits in the current branch.
type ResetCommitsInCurrentBranch struct {
	Parent gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *ResetCommitsInCurrentBranch) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.RemoveCommitsInCurrentBranch(self.Parent)
}
