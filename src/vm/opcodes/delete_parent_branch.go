package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// DeleteParentBranch removes the parent branch entry in the Git Town configuration.
type DeleteParentBranch struct {
	Branch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *DeleteParentBranch) Run(args shared.RunArgs) error {
	args.Runner.Config.RemoveParent(self.Branch)
	return nil
}
