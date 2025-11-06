package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// BranchCurrentResetToSHAIfNeeded undoes all commits on the current branch
// all the way until the given SHA.
type BranchCurrentResetToSHAIfNeeded struct {
	MustHaveSHA gitdomain.SHA
	SetToSHA    gitdomain.SHA
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
		currentBranch, err := args.Git.CurrentBranch(args.Backend)
		if err != nil {
			return err
		}
		return fmt.Errorf(messages.BranchHasWrongSHA, currentBranch.StringOr(currentSHA.Truncate(7).String()), self.SetToSHA, self.MustHaveSHA, currentSHA)
	}
	args.PrependOpcodes(&BranchCurrentResetToSHA{
		SHA: self.SetToSHA,
	})
	return nil
}
