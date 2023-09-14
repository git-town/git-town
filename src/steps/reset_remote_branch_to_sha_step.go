package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// ResetRemoteBranchToSHAStep sets the given remote branch to the given SHA,
// but only if it currently has a particular SHA.
type ResetRemoteBranchToSHAStep struct {
	Branch      domain.RemoteBranchName
	MustHaveSHA domain.SHA
	SetToSHA    domain.SHA
	EmptyStep
}

func (step *ResetRemoteBranchToSHAStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Frontend.ResetRemoteBranchToSHA(step.Branch, step.SHA)
}
