package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// ResetCurrentBranchToSHA undoes all commits on the current branch
// all the way until the given SHA.
type ResetCurrentBranchToSHA struct {
	Hard        bool
	MustHaveSHA domain.SHA
	SetToSHA    domain.SHA
	undeclaredOpcodeMethods
}

func (op *ResetCurrentBranchToSHA) Run(args shared.RunArgs) error {
	currentSHA, err := args.Runner.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	if currentSHA == op.SetToSHA {
		// nothing to do
		return nil
	}
	if currentSHA != op.MustHaveSHA {
		currentBranchName, err := args.Runner.Backend.CurrentBranch()
		if err != nil {
			return err
		}
		return fmt.Errorf(messages.BranchHasWrongSHA, currentBranchName, op.SetToSHA, op.MustHaveSHA, currentSHA)
	}
	return args.Runner.Frontend.ResetCurrentBranchToSHA(op.SetToSHA, op.Hard)
}
