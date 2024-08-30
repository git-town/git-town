package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// registers the branch with the given name as a contribution branch in the Git config
type AddToContributionBranches struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *AddToContributionBranches) Run(args shared.RunArgs) error {
	return args.Config.AddToContributionBranches(self.Branch)
}
