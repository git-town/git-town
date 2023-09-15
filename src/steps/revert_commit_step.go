package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/slice"
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
	commitsInCurrentBranch, err := args.Runner.Backend.CommitsInBranch(currentBranch, parent)
	if err != nil {
		return err
	}
	if !slice.Contains(commitsInCurrentBranch, step.SHA) {
		return fmt.Errorf(messages.BranchDoesntContainCommit, currentBranch, step.SHA, commitsInCurrentBranch.Join("|"))
	}

	// Ensure that the current branch contains the given commit?
	return args.Runner.Frontend.RevertCommit(step.SHA)
}
