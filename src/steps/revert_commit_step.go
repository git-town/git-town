package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// RevertCommitStep adds a commit to the current branch
// that reverts the commit with the given SHA.
type RevertCommitStep struct {
	SHA domain.SHA
	EmptyStep
}

func (step *RevertCommitStep) Run(args RunArgs) error {
	commits := args.Runner.Backend.Commits
	// Ensure that the current branch contains the given commit?
	return args.Runner.Frontend.RevertCommit(step.SHA)
}
