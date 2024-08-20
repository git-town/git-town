package opcodes

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// RemoveFromParkedBranches removes the branch with the given name as a parked branch.
type RemoveFromParkedBranches struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RemoveFromParkedBranches) Run(args shared.RunArgs) error {
	return args.Config.RemoveFromParkedBranches(self.Branch)
}
