package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// LineageParentSet changes the parent of the given branch to the given parent.
// Use SetParent to set the parent if no parent existed before.
type LineageParentSet struct {
	Branch gitdomain.LocalBranchName
	Parent gitdomain.LocalBranchName
}

func (self *LineageParentSet) Run(args shared.RunArgs) error {
	if err := args.Config.Value.NormalConfig.SetParent(args.Backend, self.Branch, self.Parent); err != nil {
		return err
	}
	args.FinalMessages.AddF(messages.BranchParentChanged, self.Branch, self.Parent)
	return nil
}
