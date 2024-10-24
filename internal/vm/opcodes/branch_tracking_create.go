package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// BranchTrackingCreate pushes the given local branch up to origin
// and marks it as tracking the current branch.
type BranchTrackingCreate struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchTrackingCreate) Run(args shared.RunArgs) error {
	return args.Git.CreateTrackingBranch(args.Frontend, self.Branch, gitdomain.RemoteOrigin, args.Config.Value.NormalConfig.NoPushHook())
}
