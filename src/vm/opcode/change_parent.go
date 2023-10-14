package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
)

// ChangeParent changes the parent of the given branch to the given parent.
// Use SetParent to set the parent if no parent existed before.
type ChangeParent struct {
	Branch domain.LocalBranchName
	Parent domain.LocalBranchName
	Empty
}

func (step *ChangeParent) Run(args RunArgs) error {
	err := args.Runner.Config.SetParent(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	args.Runner.FinalMessages.Add(fmt.Sprintf(messages.BranchParentChanged, step.Branch, step.Parent))
	return nil
}
