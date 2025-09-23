package opcodes

import (
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

type PushCurrentBranchForce struct {
	ForceIfIncludes bool
}

func (self *PushCurrentBranchForce) Run(args shared.RunArgs) error {
	return args.Git.ForcePushBranchSafely(args.Frontend, args.Config.Value.NormalConfig.PushHook, self.ForceIfIncludes)
}
