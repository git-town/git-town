package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// BranchLocalDelete deletes the branch with the given name.
type BranchLocalDelete struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchLocalDelete) Run(args shared.RunArgs) error {
	return args.Git.DeleteLocalBranch(args.Frontend, self.Branch)
}
