package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// ResetCurrentBranchToSHAStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetCurrentBranchToSHAStep struct {
	Hard bool
	SHA  domain.SHA
	EmptyStep
}

func (step *ResetCurrentBranchToSHAStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	currentSHA, err := run.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	if step.SHA == currentSHA {
		return nil
	}
	return run.Frontend.ResetCurrentBranchToSHA(step.SHA, step.Hard)
}
