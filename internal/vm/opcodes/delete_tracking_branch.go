package opcodes

import (
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	"github.com/git-town/git-town/v14/internal/vm/shared"
)

// DeleteTrackingBranch deletes the tracking branch of the given local branch.
type DeleteTrackingBranch struct {
	Branch                  gitdomain.RemoteBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *DeleteTrackingBranch) Run(args shared.RunArgs) error {
	return args.Git.DeleteTrackingBranch(args.Frontend, self.Branch)
}
