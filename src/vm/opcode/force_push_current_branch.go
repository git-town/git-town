package opcode

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// ForcePushCurrentBranch force-pushes the branch with the given name to the origin remote.
type ForcePushCurrentBranch struct {
	NoPushHook configdomain.NoPushHook
	undeclaredOpcodeMethods
}

func (self *ForcePushCurrentBranch) Run(args shared.RunArgs) error {
	currentBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	shouldPush, err := args.Runner.Backend.ShouldPushBranch(currentBranch, currentBranch.TrackingBranch())
	if err != nil {
		return err
	}
	if !shouldPush && !args.Runner.GitTown.DryRun {
		return nil
	}
	return args.Runner.Frontend.ForcePushBranch(self.NoPushHook)
}
