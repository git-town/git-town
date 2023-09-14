package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// CreateRemoteBranchStep pushes the given local branch up to origin.
type CreateRemoteBranchStep struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	SHA        domain.SHA
	EmptyStep
}

func (step *CreateRemoteBranchStep) Run(args RunArgs) error {
	return args.Run.Frontend.CreateRemoteBranch(step.SHA, step.Branch, step.NoPushHook)
}
