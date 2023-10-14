package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// FetchUpstream brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type FetchUpstream struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (step *FetchUpstream) Run(args RunArgs) error {
	return args.Runner.Frontend.FetchUpstream(step.Branch)
}
