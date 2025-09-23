package opcodes

import (
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// PushCurrentBranch pushes the current branch to its existing tracking branch.
type PushCurrentBranch struct{}

func (self *PushCurrentBranch) Run(args shared.RunArgs) error {
	return args.Git.PushCurrentBranch(args.Frontend, args.Config.Value.NormalConfig.PushHook)
}
