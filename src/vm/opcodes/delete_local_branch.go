package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// DeleteLocalBranch deletes the branch with the given name.
type DeleteLocalBranch struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *DeleteLocalBranch) Run(args shared.RunArgs) error {
	return args.Git.DeleteLocalBranch(args.Frontend, self.Branch)
}
