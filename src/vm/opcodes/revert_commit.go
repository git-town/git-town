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
	SHA gitdomain.SHA
	undeclaredOpcodeMethods
}

func (self *RevertCommit) Run(args shared.RunArgs) error {
	currentBranch, err := args.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	parent := args.Lineage.Parent(currentBranch)
	commitsInCurrentBranch, err := args.Backend.CommitsInBranch(currentBranch, parent)
	if err != nil {
		return err
	}
	if !commitsInCurrentBranch.ContainsSHA(self.SHA) {
		return fmt.Errorf(messages.BranchDoesntContainCommit, currentBranch, self.SHA, commitsInCurrentBranch.SHAs().Join("|"))
	}
	return args.Frontend.RevertCommit(self.SHA)
}
