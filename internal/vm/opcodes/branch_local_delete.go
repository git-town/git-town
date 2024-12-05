package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// BranchLocalDelete deletes the branch with the given name.
type BranchLocalDelete struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchLocalDelete) Run(args shared.RunArgs) error {
	args.FinalMessages.Add(fmt.Sprintf(messages.BranchDeleted, self.Branch))
	return args.Git.DeleteLocalBranch(args.Frontend, self.Branch)
}
