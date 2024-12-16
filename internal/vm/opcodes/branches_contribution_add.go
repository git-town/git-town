package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// registers the branch with the given name as a contribution branch in the Git config
type BranchesContributionAdd struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchesContributionAdd) Run(args shared.RunArgs) error {
	return args.Config.Value.NormalConfig.AddToContributionBranches(self.Branch)
}
