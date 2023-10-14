package step

import "github.com/git-town/git-town/v9/src/domain"

// SetParentIfBranchExists sets the given parent branch as the parent of the given branch,
// but only the latter exists. This is useful when the branch might or might not be deleted
// depending on conditions evaluated at runtime.
// Use SetParent if you are sure that the branch will exist.
// Use ChangeParent to change an existing parent.
type SetParentIfBranchExists struct {
	Branch domain.LocalBranchName
	Parent domain.LocalBranchName
	Empty
}

func (step *SetParentIfBranchExists) Run(args RunArgs) error {
	branchExists := args.Runner.Backend.BranchExists(step.Branch)
	if !branchExists {
		return nil
	}
	return args.Runner.Config.SetParent(step.Branch, step.Parent)
}
