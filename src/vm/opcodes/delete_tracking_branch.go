package opcodes

import (
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/vm/shared"
)

// DeleteTrackingBranch deletes the tracking branch of the given local branch.
type DeleteTrackingBranch struct {
	Branch gitdomain.RemoteBranchName
	undeclaredOpcodeMethods
}

func (self *DeleteTrackingBranch) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.DeleteTrackingBranch(self.Branch)
}
