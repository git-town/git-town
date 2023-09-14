package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// ResetRemoteBranchToSHAStep sets the given remote branch to the given SHA,
// but only if it currently has a particular SHA.
type ResetRemoteBranchToSHAStep struct {
	Branch      domain.RemoteBranchName
	MustHaveSHA domain.SHA
	SetToSHA    domain.SHA
	EmptyStep
}

func (step *ResetRemoteBranchToSHAStep) Run(args RunArgs) error {
	return args.Runner.Frontend.ResetRemoteBranchToSHA(step.Branch, step.SetToSHA)
}
