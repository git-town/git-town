package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// CreateBranchStep creates a new branch but leaves the current branch unchanged.
type CreateBranchStep struct {
	Branch        domain.LocalBranchName
	StartingPoint domain.Location
	EmptyStep
}

func (step *CreateBranchStep) Run(args RunArgs) error {
	return args.Runner.Frontend.CreateBranch(step.Branch, step.StartingPoint)
}
