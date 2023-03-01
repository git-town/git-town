package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// FetchUpstreamStep brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type FetchUpstreamStep struct {
	EmptyStep
	Branch string
}

func (step *FetchUpstreamStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return repo.Logging.FetchUpstream(step.Branch)
}
