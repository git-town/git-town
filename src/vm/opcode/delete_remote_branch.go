package opcode

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// DeleteRemoteBranch deletes the tracking branch of the given local branch.
type DeleteRemoteBranch struct {
	Branch gitdomain.RemoteBranchName
	undeclaredOpcodeMethods
}

func (self *DeleteRemoteBranch) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.DeleteRemoteBranch(self.Branch)
}
