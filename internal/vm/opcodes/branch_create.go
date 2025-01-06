package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// BranchCreate creates a new branch but leaves the current branch unchanged.
type BranchCreate struct {
	Branch                  gitdomain.LocalBranchName
	StartingPoint           gitdomain.Location
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchCreate) Run(args shared.RunArgs) error {
	return args.Git.CreateBranch(args.Frontend, self.Branch, self.StartingPoint)
}
