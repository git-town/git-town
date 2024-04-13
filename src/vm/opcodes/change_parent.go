package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// ChangeParent changes the parent of the given branch to the given parent.
// Use SetParent to set the parent if no parent existed before.
type ChangeParent struct {
	Branch gitdomain.LocalBranchName
	Parent gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *ChangeParent) Run(args shared.RunArgs) error {
	err := args.Runner.Config.SetParent(self.Branch, self.Parent)
	if err != nil {
		return err
	}
	args.Runner.FinalMessages.Add(fmt.Sprintf(messages.BranchParentChanged, self.Branch, self.Parent))
	return nil
}
