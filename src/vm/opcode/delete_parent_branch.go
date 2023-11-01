package opcode

import (
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/vm/shared"
)

// DeleteParentBranch removes the parent branch entry in the Git Town configuration.
type DeleteParentBranch struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *DeleteParentBranch) Run(args shared.RunArgs) error {
	args.Runner.Config.RemoveParent(self.Branch)
	return nil
}
