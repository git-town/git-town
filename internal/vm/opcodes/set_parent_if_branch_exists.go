package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// SetParentIfBranchExists sets the given parent branch as the parent of the given branch,
// but only the latter exists.
type SetParentIfBranchExists struct {
	Branch                  gitdomain.LocalBranchName
	Parent                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SetParentIfBranchExists) Run(args shared.RunArgs) error {
	if !args.Git.BranchExists(args.Backend, self.Branch) {
		return nil
	}
	args.PrependOpcodes(&SetParent{Branch: self.Branch, Parent: self.Parent})
	return nil
}
