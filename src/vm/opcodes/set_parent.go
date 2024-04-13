package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// SetParent sets the given parent branch as the parent of the given branch.
// Use ChangeParent to change an existing parent.
type SetParent struct {
	Branch gitdomain.LocalBranchName
	Parent gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *SetParent) Run(args shared.RunArgs) error {
	return args.Runner.Config.SetParent(self.Branch, self.Parent)
}
