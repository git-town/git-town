package steps

import (
	"github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/src/hosting"
)

// ContinueMergeStep finishes an ongoing merge conflict
// assuming all conflicts have been resolved by the user.
type ContinueMergeStep struct {
	EmptyStep
}

func (step *ContinueMergeStep) CreateContinueStep() Step {
	return step
}

func (step *ContinueMergeStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	if run.Backend.HasMergeInProgress() {
		return run.Frontend.CommitNoEdit()
	}
	return nil
}
