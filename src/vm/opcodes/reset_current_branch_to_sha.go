package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// ResetCurrentBranchToSHA undoes all commits on the current branch
// all the way until the given SHA.
type ResetCurrentBranchToSHA struct {
	Hard                    bool
	MustHaveSHA             gitdomain.SHA
	SetToSHA                gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ResetCurrentBranchToSHA) Run(args shared.RunArgs) error {
	currentSHA, err := args.Git.CurrentSHA(args.Backend)
	if err != nil {
		return err
	}
	if currentSHA == self.SetToSHA {
		// nothing to do
		return nil
	}
	if currentSHA != self.MustHaveSHA {
		currentBranchName, err := args.Git.CurrentBranch(args.Backend)
		if err != nil {
			return err
		}
		return fmt.Errorf(messages.BranchHasWrongSHA, currentBranchName, self.SetToSHA, self.MustHaveSHA, currentSHA)
	}
	return args.Git.ResetCurrentBranchToSHA(args.Frontend, self.SetToSHA, self.Hard)
}
