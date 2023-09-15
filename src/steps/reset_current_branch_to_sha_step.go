package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
)

// ResetCurrentBranchToSHAStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetCurrentBranchToSHAStep struct {
	Hard        bool
	MustHaveSHA domain.SHA
	SetToSHA    domain.SHA
	EmptyStep
}

func (step *ResetCurrentBranchToSHAStep) Run(args RunArgs) error {
	currentSHA, err := args.Runner.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	if step.SetToSHA == currentSHA {
		// nothing to do here
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
