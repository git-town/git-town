package opcodes

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// RemoveFromContributionBranches removes the branch with the given name as a contribution branch.
type RemoveFromContributionBranches struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RemoveFromContributionBranches) Run(args shared.RunArgs) error {
	return args.Config.RemoveFromContributionBranches(self.Branch)
}
