package opcodes

import (
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	"github.com/git-town/git-town/v14/internal/vm/shared"
)

// CreateBranch creates a new branch but leaves the current branch unchanged.
type CreateBranch struct {
	Branch                  gitdomain.LocalBranchName
	StartingPoint           gitdomain.Location
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CreateBranch) Run(args shared.RunArgs) error {
	return args.Git.CreateBranch(args.Frontend, self.Branch, self.StartingPoint)
}
