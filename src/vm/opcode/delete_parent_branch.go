package opcode

import (
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// DeleteParentBranch removes the parent branch entry in the Git Town configuration.
type DeleteParentBranch struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *DeleteParentBranch) Run(args shared.RunArgs) error {
	args.Runner.GitTown.RemoveParent(self.Branch)
	return nil
}
