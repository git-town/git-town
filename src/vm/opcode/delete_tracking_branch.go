package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// DeleteTrackingBranch deletes the tracking branch of the given local branch.
type DeleteTrackingBranch struct {
	Branch domain.RemoteBranchName
	undeclaredOpcodeMethods
}

func (step *DeleteTrackingBranch) Run(args RunArgs) error {
	return args.Runner.Frontend.DeleteRemoteBranch(step.Branch)
}
