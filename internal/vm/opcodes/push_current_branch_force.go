package opcodes

import (
	"github.com/git-town/git-town/v18/internal/vm/shared"
)

type PushCurrentBranchForce struct {
	ForceIfIncludes         bool
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *PushCurrentBranchForce) Run(args shared.RunArgs) error {
	return args.Git.ForcePushBranchSafely(args.Frontend, args.Config.Value.NormalConfig.NoPushHook(), self.ForceIfIncludes)
}
