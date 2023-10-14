package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
)

// ResetCurrentBranchToSHA undoes all commits on the current branch
// all the way until the given SHA.
type ResetCurrentBranchToSHA struct {
	Hard        bool
	MustHaveSHA domain.SHA
	SetToSHA    domain.SHA
	BaseOpcode
}

func (step *ResetCurrentBranchToSHA) Run(args RunArgs) error {
	currentSHA, err := args.Runner.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	if currentSHA == step.SetToSHA {
		// nothing to do
		return nil
	}
	if currentSHA != step.MustHaveSHA {
		currentBranchName, err := args.Runner.Backend.CurrentBranch()
		if err != nil {
			return err
		}
		return fmt.Errorf(messages.BranchHasWrongSHA, currentBranchName, step.SetToSHA, step.MustHaveSHA, currentSHA)
	}
	return args.Runner.Frontend.ResetCurrentBranchToSHA(step.SetToSHA, step.Hard)
}
