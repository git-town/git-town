package steps

import (
	"github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/src/hosting"
)

// ContinueRebaseStep finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebaseStep struct {
	EmptyStep
}

func (step *ContinueRebaseStep) CreateAbortStep() Step {
	return &AbortRebaseStep{}
}

func (step *ContinueRebaseStep) CreateContinueStep() Step {
	return step
}

func (step *ContinueRebaseStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	hasRebaseInProgress, err := run.Backend.HasRebaseInProgress()
	if err != nil {
		return err
	}
	if hasRebaseInProgress {
		return run.Frontend.ContinueRebase()
	}
	return nil
}
