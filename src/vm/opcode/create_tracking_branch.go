package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// CreateTrackingBranch pushes the given local branch up to origin
// and marks it as tracking the current branch.
type CreateTrackingBranch struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	undeclaredOpcodeMethods
}

func (step *CreateTrackingBranch) Run(args RunArgs) error {
	return args.Runner.Frontend.CreateTrackingBranch(step.Branch, domain.OriginRemote, step.NoPushHook)
}
