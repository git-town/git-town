package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// BranchCurrentResetToSHAIfNeeded undoes all commits on the current branch
// all the way until the given SHA.
type BranchCurrentResetToSHAIfNeeded struct {
	MustHaveSHA             gitdomain.SHA
	SetToSHA                gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchCurrentResetToSHAIfNeeded) Run(args shared.RunArgs) error {
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
	args.PrependOpcodes(&BranchCurrentResetToSHA{
		SetToSHA: self.SetToSHA,
	})
	return nil
}
