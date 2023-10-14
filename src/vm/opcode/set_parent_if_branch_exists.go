package opcode

import "github.com/git-town/git-town/v9/src/domain"

// SetParentIfBranchExists sets the given parent branch as the parent of the given branch,
// but only the latter exists.
type SetParentIfBranchExists struct {
	Branch domain.LocalBranchName
	Parent domain.LocalBranchName
	Empty
}

func (step *SetParentIfBranchExists) Run(args RunArgs) error {
	if !args.Runner.Backend.BranchExists(step.Branch) {
		return nil
	}
	return args.Runner.Config.SetParent(step.Branch, step.Parent)
}
