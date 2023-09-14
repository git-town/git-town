package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// RevertCommitStep reverts the commit with the given sha.
type RevertCommitStep struct {
	SHA domain.SHA
	EmptyStep
}

func (step *RevertCommitStep) Run(args RunArgs) error {
	return args.Runner.Frontend.RevertCommit(step.SHA)
}
