package opcode

import (
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// DeleteRemoteBranch deletes the tracking branch of the given local branch.
type DeleteRemoteBranch struct {
	Branch domain.RemoteBranchName
	undeclaredOpcodeMethods
}

func (self *DeleteRemoteBranch) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.DeleteRemoteBranch(self.Branch)
}
