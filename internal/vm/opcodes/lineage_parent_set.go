package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// LineageParentSet changes the parent of the given branch to the given parent.
// Use SetParent to set the parent if no parent existed before.
type LineageParentSet struct {
	Branch                  gitdomain.LocalBranchName
	Parent                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *LineageParentSet) Run(args shared.RunArgs) error {
	if err := args.Config.Value.NormalConfig.SetParent(args.Backend, self.Branch, self.Parent); err != nil {
		return err
	}
	args.FinalMessages.Add(fmt.Sprintf(messages.BranchParentChanged, self.Branch, self.Parent))
	return nil
}
