package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// PushCurrentBranchIfLocal pushes the current branch to its existing tracking branch.
type PushCurrentBranchIfLocal struct {
	CurrentBranch gitdomain.LocalBranchName
}

func (self *PushCurrentBranchIfLocal) Run(args shared.RunArgs) error {
	hasTrackingBranch := args.Git.CurrentBranchHasTrackingBranch(args.Backend)
	if !hasTrackingBranch {
		args.PrependOpcodes(&BranchTrackingCreate{
			Branch: self.CurrentBranch,
		})
	}
	return nil
}
