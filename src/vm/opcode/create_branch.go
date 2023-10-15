package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// CreateBranch creates a new branch but leaves the current branch unchanged.
type CreateBranch struct {
	Branch        domain.LocalBranchName
	StartingPoint domain.Location
	undeclaredOpcodeMethods
}

func (op *CreateBranch) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.CreateBranch(op.Branch, op.StartingPoint)
}
