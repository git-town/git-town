package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// SetParent sets the given parent branch as the parent of the given branch.
// Use ChangeParent to change an existing parent.
type SetParent struct {
	Branch domain.LocalBranchName
	Parent domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (op *SetParent) Run(args shared.RunArgs) error {
	return args.Runner.Config.SetParent(op.Branch, op.Parent)
}
