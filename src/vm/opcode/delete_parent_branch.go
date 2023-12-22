package opcode

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// DeleteParentBranch removes the parent branch entry in the Git Town configuration.
type DeleteParentBranch struct {
	Branch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *DeleteParentBranch) Run(args shared.RunArgs) error {
	args.Runner.GitTown.RemoveParent(self.Branch)
	return nil
}
