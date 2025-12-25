package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// BranchTrackingCreateIfNeeded creates the tracking branch for the current branch if it doesn't exist.
type BranchTrackingCreateIfNeeded struct {
	CurrentBranch gitdomain.LocalBranchName
}

func (self *BranchTrackingCreateIfNeeded) Run(args shared.RunArgs) error {
	hasTrackingBranch := args.Git.CurrentBranchHasTrackingBranch(args.Backend)
	if !hasTrackingBranch {
		args.PrependOpcodes(&BranchTrackingCreate{
			Branch: self.CurrentBranch,
		})
	}
	return nil
}
