package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// BranchLocalDelete deletes the branch with the given name.
type BranchLocalDelete struct {
	Branch gitdomain.LocalBranchName
}

func (self *BranchLocalDelete) Run(args shared.RunArgs) error {
	args.FinalMessages.AddF(messages.BranchDeleted, self.Branch)
	return args.Git.DeleteLocalBranch(args.Frontend, self.Branch)
}
