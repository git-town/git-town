package steps

import (
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// FetchUpstreamStep brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type FetchUpstreamStep struct {
	NoOpStep
	BranchName string
}

// Run executes this step.
func (step *FetchUpstreamStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	return repo.Logging.FetchUpstream(step.BranchName)
}
