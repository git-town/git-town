package opcodes

import (
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// ForcePushCurrentBranch force-pushes the branch with the given name to the origin remote.
// TODO: rename to PushCurrentBranchForce
type ForcePush struct {
	ForceIfIncludes         bool
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ForcePush) Run(args shared.RunArgs) error {
	return args.Git.ForcePushBranchSafely(args.Frontend, args.Config.Value.NormalConfig.NoPushHook(), self.ForceIfIncludes)
}
