package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// SetParentIfBranchExists sets the given parent branch as the parent of the given branch,
// but only the latter exists.
type SetParentIfBranchExists struct {
	Branch                  gitdomain.LocalBranchName
	Parent                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SetParentIfBranchExists) Run(args shared.RunArgs) error {
	if !args.Backend.BranchExists(self.Branch) {
		return nil
	}
	return args.Config.SetParent(self.Branch, self.Parent)
}
