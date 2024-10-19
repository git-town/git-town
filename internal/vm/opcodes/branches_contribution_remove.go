package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// removes the branch with the given name from the contribution branches list in the Git config
type BranchesContributionRemove struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchesContributionRemove) Run(args shared.RunArgs) error {
	var err error
	if args.Config.NormalConfig.ContributionBranches.Contains(self.Branch) {
		err = args.Config.NormalConfig.RemoveFromContributionBranches(self.Branch)
	}
	return err
}
