package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/messages"
	"github.com/git-town/git-town/v13/src/vm/shared"
)

// ResetCurrentBranchToSHA undoes all commits on the current branch
// all the way until the given SHA.
type ResetCurrentBranchToSHA struct {
	Hard        bool
	MustHaveSHA gitdomain.SHA
	SetToSHA    gitdomain.SHA
	undeclaredOpcodeMethods
}

func (self *ResetCurrentBranchToSHA) Run(args shared.RunArgs) error {
	currentSHA, err := args.Runner.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	if currentSHA == self.SetToSHA {
		// nothing to do
		return nil
	}
	if currentSHA != self.MustHaveSHA {
		currentBranchName, err := args.Runner.Backend.CurrentBranch()
		if err != nil {
			return err
		}
		return fmt.Errorf(messages.BranchHasWrongSHA, currentBranchName, self.SetToSHA, self.MustHaveSHA, currentSHA)
	}
	return args.Runner.Frontend.ResetCurrentBranchToSHA(self.SetToSHA, self.Hard)
}
