package opcodes

import (
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/vm/shared"
)

// BranchLocalRename renames a local branch.
type BranchLocalRename struct {
	NewName gitdomain.LocalBranchName
	OldName gitdomain.LocalBranchName
}

func (self *BranchLocalRename) Run(args shared.RunArgs) error {
	return args.Git.RenameBranch(args.Frontend, self.OldName, self.NewName)
}
