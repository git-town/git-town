package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// DeleteBranchIfEmptyAtRuntime allows running different opcodes based on a condition evaluated at runtime.
type DeleteBranchIfEmptyAtRuntime struct {
	Branch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *DeleteBranchIfEmptyAtRuntime) Run(args shared.RunArgs) error {
	parentPtr := args.Lineage.Parent(self.Branch)
	if parentPtr != nil {
		parent := *parentPtr
		hasUnmergedChanges, err := args.Runner.Backend.BranchHasUnmergedChanges(self.Branch, parent)
		if err != nil {
			return err
		}
		if hasUnmergedChanges {
			args.PrependOpcodes(&QueueMessage{
				Message: fmt.Sprintf(messages.BranchDeletedHasUnmergedChanges, self.Branch),
			})
			return nil
		}
	}
	args.PrependOpcodes(
		&CheckoutParent{CurrentBranch: self.Branch},
		&DeleteLocalBranch{Branch: self.Branch},
		&RemoveBranchFromLineage{
			Branch: self.Branch,
		},
		&QueueMessage{
			Message: fmt.Sprintf(messages.BranchDeleted, self.Branch),
		},
	)
	return nil
}
