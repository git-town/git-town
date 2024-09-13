package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// CreateBranch creates a new branch but leaves the current branch unchanged.
type RenameBranch struct {
	NewName                 gitdomain.LocalBranchName
	OldName                 gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RenameBranch) Run(args shared.RunArgs) error {
	return args.Git.RenameBranch(args.Frontend, self.OldName, self.NewName)
}
