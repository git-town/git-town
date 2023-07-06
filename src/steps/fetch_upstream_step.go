package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// FetchUpstreamStep brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type FetchUpstreamStep struct {
	EmptyStep
	Branch string
}

func (step *FetchUpstreamStep) Run(run *git.ProdRunner, _connector hosting.Connector) error {
	return run.Frontend.FetchUpstream(step.Branch)
}
