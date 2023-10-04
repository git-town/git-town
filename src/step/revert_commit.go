package step

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/git-town/git-town/v9/src/messages"
)

// RevertCommit adds a commit to the current branch
// that reverts the commit with the given SHA.
type RevertCommit struct {
	SHA domain.SHA
	Empty
}

func (step *RevertCommit) Run(args RunArgs) error {
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
	return args.Runner.Frontend.RevertCommit(step.SHA)
}
