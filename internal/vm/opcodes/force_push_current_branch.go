package opcodes

import (
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// ForcePushCurrentBranch force-pushes the branch with the given name to the origin remote.
type ForcePushCurrentBranch struct {
	ForceIfIncludes         bool
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ForcePushCurrentBranch) Run(args shared.RunArgs) error {
	currentBranch, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	shouldPush, err := args.Git.ShouldPushBranch(args.Backend, currentBranch)
	if err != nil {
		return err
	}
	if !shouldPush {
		return nil
	}
	return args.Git.ForcePushBranchSafely(args.Frontend, args.Config.Config.NoPushHook(), self.ForceIfIncludes)
}
