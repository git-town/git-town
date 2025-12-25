package opcodes

import (
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// PushCurrentBranchForceIgnoreError attempts to force-push.
type PushCurrentBranchForceIgnoreError struct{}

func (self *PushCurrentBranchForceIgnoreError) Run(args shared.RunArgs) error {
	_ = args.Git.ForcePushBranchSafely(args.Frontend, args.Config.Value.NormalConfig.PushHook, true)
	return nil
}
