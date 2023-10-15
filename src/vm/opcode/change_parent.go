package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// ChangeParent changes the parent of the given branch to the given parent.
// Use SetParent to set the parent if no parent existed before.
type ChangeParent struct {
	Branch domain.LocalBranchName
	Parent domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (op *ChangeParent) Run(args shared.RunArgs) error {
	err := args.Runner.Config.SetParent(op.Branch, op.Parent)
	if err != nil {
		return err
	}
	args.Runner.FinalMessages.Add(fmt.Sprintf(messages.BranchParentChanged, op.Branch, op.Parent))
	return nil
}
