package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// DeleteRemoteBranch deletes the tracking branch of the given local branch.
type DeleteRemoteBranch struct {
	Branch domain.RemoteBranchName
	Empty
}

func (step *DeleteRemoteBranch) Run(args RunArgs) error {
	return args.Runner.Frontend.DeleteRemoteBranch(step.Branch)
}
