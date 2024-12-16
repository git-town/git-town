package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// BranchTrackingCreate pushes the given local branch up to origin
// and marks it as tracking the current branch.
type BranchTrackingCreate struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchTrackingCreate) Run(args shared.RunArgs) error {
	return args.Git.CreateTrackingBranch(args.Frontend, self.Branch, args.Config.Value.NormalConfig.DevRemote, args.Config.Value.NormalConfig.NoPushHook())
}
