package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// CreateRemoteBranch pushes the given local branch up to origin.
type CreateRemoteBranch struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	SHA        domain.SHA
	BaseOpcode
}

func (step *CreateRemoteBranch) Run(args RunArgs) error {
	return args.Runner.Frontend.CreateRemoteBranch(step.SHA, step.Branch, step.NoPushHook)
}
