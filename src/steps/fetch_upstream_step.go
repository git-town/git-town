package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// FetchUpstreamStep brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type FetchUpstreamStep struct {
	Branch domain.LocalBranchName
	EmptyStep
}

func (step *FetchUpstreamStep) Run(args RunArgs) error {
	return args.Run.Frontend.FetchUpstream(step.Branch)
}
