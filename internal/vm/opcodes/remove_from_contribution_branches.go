package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// removes the branch with the given name from the contribution branches list in the Git config
type RemoveFromContributionBranches struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RemoveFromContributionBranches) Run(args shared.RunArgs) error {
	return args.Config.RemoveFromContributionBranches(self.Branch)
}
