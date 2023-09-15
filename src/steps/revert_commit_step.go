package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
)

// RevertCommitStep adds a commit to the current branch
// that reverts the commit with the given SHA.
type RevertCommitStep struct {
	SHA domain.SHA
	EmptyStep
}

func (step *RevertCommitStep) Run(args RunArgs) error {
	currentBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	parent := args.Lineage.Parent(currentBranch)
	commitsInCurrentBranch := args.Runner.Backend.CommitsInCurrentBranch(currentBranch, parent)
	if !commitsInCurrentBranch.Contains(step.SHA) {
		return fmt.Errorf("branch %q does not contain commit %q. Found commits %s", currentBranch, step.SHA)
	}

	// Ensure that the current branch contains the given commit?
	return args.Runner.Frontend.RevertCommit(step.SHA)
}
