package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// DeleteTrackingBranch deletes the tracking branch of the given local branch.
type DeleteTrackingBranch struct {
	Branch                  gitdomain.RemoteBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *DeleteTrackingBranch) Run(args shared.RunArgs) error {
	_ = args.Git.DeleteTrackingBranch(args.Frontend, self.Branch)
	return nil
}
