package opcode

import (
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/vm/shared"
)

// CreateTrackingBranch pushes the given local branch up to origin
// and marks it as tracking the current branch.
type CreateTrackingBranch struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	undeclaredOpcodeMethods
}

func (self *CreateTrackingBranch) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.CreateTrackingBranch(self.Branch, domain.OriginRemote, self.NoPushHook)
}
