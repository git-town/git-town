package opcodes

import (
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// ForcePushCurrentBranch force-pushes the branch with the given name to the origin remote.
type ForcePushCurrentBranch struct {
	undeclaredOpcodeMethods
}

func (self *ForcePushCurrentBranch) Run(args shared.RunArgs) error {
	currentBranch, err := args.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	shouldPush, err := args.Backend.ShouldPushBranch(currentBranch, currentBranch.TrackingBranch())
	if err != nil {
		return err
	}
	if !shouldPush {
		return nil
	}
	return args.Frontend.ForcePushBranchSafely(args.Config.Config.NoPushHook())
}
