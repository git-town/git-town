package opcodes

import (
	"github.com/git-town/git-town/v19/internal/vm/shared"
)

// PushCurrentBranchForceIgnoreError attempts to force-push.
type PushCurrentBranchForceIgnoreError struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *PushCurrentBranchForceIgnoreError) Run(args shared.RunArgs) error {
	_ = args.Git.ForcePushBranchSafely(args.Frontend, args.Config.Value.NormalConfig.NoPushHook(), true)
	return nil
}
