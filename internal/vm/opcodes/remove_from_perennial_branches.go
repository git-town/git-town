package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// removes the branch with the given name from the perennial branches list in the Git config
type RemoveFromPerennialBranches struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RemoveFromPerennialBranches) Run(args shared.RunArgs) error {
	var err error
	if args.Config.Config.PerennialBranches.Contains(self.Branch) {
		err = args.Config.RemoveFromPerennialBranches(self.Branch)
	}
	return err
}
