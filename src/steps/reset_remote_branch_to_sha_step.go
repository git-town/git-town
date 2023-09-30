package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
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
	currentSHA, err := args.Runner.Backend.SHAForBranch(step.Branch.BranchName())
	if err != nil {
		return err
	}
	if currentSHA != step.MustHaveSHA {
		return fmt.Errorf(messages.BranchHasWrongSHA, step.Branch, step.SetToSHA, step.MustHaveSHA, currentSHA)
	}
	return args.Runner.Frontend.ResetRemoteBranchToSHA(step.Branch, step.SetToSHA)
}
