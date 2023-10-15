package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// RevertCommit adds a commit to the current branch
// that reverts the commit with the given SHA.
type RevertCommit struct {
	SHA domain.SHA
	undeclaredOpcodeMethods
}

func (op *RevertCommit) Run(args shared.RunArgs) error {
	currentBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	parent := args.Lineage.Parent(currentBranch)
	commitsInCurrentBranch, err := args.Runner.Backend.CommitsInBranch(currentBranch, parent)
	if err != nil {
		return err
	}
	if !slice.Contains(commitsInCurrentBranch, op.SHA) {
		return fmt.Errorf(messages.BranchDoesntContainCommit, currentBranch, op.SHA, commitsInCurrentBranch.Join("|"))
	}
	return args.Runner.Frontend.RevertCommit(op.SHA)
}
