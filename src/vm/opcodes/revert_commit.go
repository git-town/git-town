package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// RevertCommit adds a commit to the current branch
// that reverts the commit with the given SHA.
type RevertCommit struct {
	SHA                     gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RevertCommit) Run(args shared.RunArgs) error {
	currentBranch, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	parent := args.Config.Config.Lineage.Parent(currentBranch)
	commitsInCurrentBranch, err := args.Git.CommitsInBranch(args.Backend, currentBranch, parent)
	if err != nil {
		return err
	}
	if !commitsInCurrentBranch.ContainsSHA(self.SHA) {
		return fmt.Errorf(messages.BranchDoesntContainCommit, currentBranch, self.SHA, commitsInCurrentBranch.SHAs().Join("|"))
	}
	return args.Git.RevertCommit(args.Frontend, self.SHA)
}
