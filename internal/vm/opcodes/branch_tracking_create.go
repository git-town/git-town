package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// BranchTrackingCreate pushes the given local branch up to origin
// and marks it as tracking the current branch.
type BranchTrackingCreate struct {
	Branch gitdomain.LocalBranchName
}

func (self *BranchTrackingCreate) Run(args shared.RunArgs) error {
	return args.Git.CreateTrackingBranch(args.Frontend, self.Branch, args.Config.Value.NormalConfig.DevRemote, args.Config.Value.NormalConfig.PushHook)
}
