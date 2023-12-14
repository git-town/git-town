package opcode

import (
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// SetParentIfBranchExists sets the given parent branch as the parent of the given branch,
// but only the latter exists.
type SetParentIfBranchExists struct {
	Branch domain.LocalBranchName
	Parent domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *SetParentIfBranchExists) Run(args shared.RunArgs) error {
	if !args.Runner.Backend.BranchExists(self.Branch) {
		return nil
	}
	return args.Runner.GitTown.SetParent(self.Branch, self.Parent)
}
