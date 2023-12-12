package opcode

import (
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// CreateBranch creates a new branch but leaves the current branch unchanged.
type CreateBranch struct {
	Branch        domain.LocalBranchName
	StartingPoint domain.Location
	undeclaredOpcodeMethods
}

func (self *CreateBranch) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.CreateBranch(self.Branch, self.StartingPoint)
}
