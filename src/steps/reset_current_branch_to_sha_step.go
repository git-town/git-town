package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
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
		return nil
	}
	return args.Runner.Frontend.ResetCurrentBranchToSHA(step.SetToSHA, step.Hard)
}
