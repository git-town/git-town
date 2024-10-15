package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// CreateBranch creates a new branch but leaves the current branch unchanged.
type BranchRename struct {
	NewName                 gitdomain.LocalBranchName
	OldName                 gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchRename) Run(args shared.RunArgs) error {
	return args.Git.Rename(args.Frontend, self.OldName, self.NewName)
}
