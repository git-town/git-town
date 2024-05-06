package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// DeleteBranchIfEmptyAtRuntime allows running different opcodes based on a condition evaluated at runtime.
type DeleteBranchIfEmptyAtRuntime struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *DeleteBranchIfEmptyAtRuntime) Run(args shared.RunArgs) error {
	parent, hasParent := args.Lineage.Parent(self.Branch).Get()
	if !hasParent {
		return nil
	}
	hasUnmergedChanges, err := args.Backend.BranchHasUnmergedChanges(self.Branch, parent)
	if err != nil {
		return err
	}
	if hasUnmergedChanges {
		args.PrependOpcodes(&QueueMessage{
			Message: fmt.Sprintf(messages.BranchDeletedHasUnmergedChanges, self.Branch),
		})
	} else {
		args.PrependOpcodes(
			&CheckoutParent{CurrentBranch: self.Branch},
			&DeleteLocalBranch{Branch: self.Branch},
			&RemoveBranchFromLineage{
				Branch: self.Branch,
			},
			&QueueMessage{
				Message: fmt.Sprintf(messages.BranchDeleted, self.Branch),
			})
	}
	return nil
}
