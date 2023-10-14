package opcode

import "github.com/git-town/git-town/v9/src/domain"

// SetParent sets the given parent branch as the parent of the given branch.
// Use ChangeParent to change an existing parent.
type SetParent struct {
	Branch domain.LocalBranchName
	Parent domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (step *SetParent) Run(args RunArgs) error {
	return args.Runner.Config.SetParent(step.Branch, step.Parent)
}
