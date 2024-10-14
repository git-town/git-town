package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// removes the branch with the given name from the parked branches list in the Git config
type RemoveFromParkedBranches struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RemoveFromParkedBranches) Run(args shared.RunArgs) error {
	var err error
	if args.Config.Config.ParkedBranches.Contains(self.Branch) {
		err = args.Config.RemoveFromParkedBranches(self.Branch)
	}
	return err
}
