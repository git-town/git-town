package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// ResetToSHAStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetToSHAStep struct {
	Hard bool
	SHA  domain.SHA
	EmptyStep
}

func (step *ResetToSHAStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	currentSHA, err := run.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	if step.SHA == currentSHA {
		return nil
	}
	return run.Frontend.ResetToSHA(step.SHA, step.Hard)
}
