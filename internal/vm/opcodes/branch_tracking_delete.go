package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// BranchTrackingDelete deletes the tracking branch of the given local branch.
type BranchTrackingDelete struct {
	Branch                  gitdomain.RemoteBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchTrackingDelete) Run(args shared.RunArgs) error {
	_ = args.Git.DeleteTrackingBranch(args.Frontend, self.Branch)
	return nil
}
