package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// RevertCommitStep reverts the commit with the given sha.
type RevertCommitStep struct {
	SHA domain.SHA
	EmptyStep
}

func (step *RevertCommitStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Frontend.RevertCommit(step.SHA)
}
