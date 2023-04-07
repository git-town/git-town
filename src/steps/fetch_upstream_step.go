package steps

import (
	"github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/src/hosting"
)

// FetchUpstreamStep brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type FetchUpstreamStep struct {
	EmptyStep
	Branch string
}

// TODO: is this used?
func (step *FetchUpstreamStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	return run.Frontend.FetchUpstream(step.Branch)
}
