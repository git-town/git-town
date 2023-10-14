package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// CreateBranch creates a new branch but leaves the current branch unchanged.
type CreateBranch struct {
	Branch        domain.LocalBranchName
	StartingPoint domain.Location
	undeclaredOpcodeMethods
}

func (step *CreateBranch) Run(args RunArgs) error {
	return args.Runner.Frontend.CreateBranch(step.Branch, step.StartingPoint)
}
