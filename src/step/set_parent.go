package step

import "github.com/git-town/git-town/v9/src/domain"

// SetParent sets the given parent branch as the parent of the given branch.
// Use ChangeParent to change an existing parent.
type SetParent struct {
	Branch domain.LocalBranchName
	Parent domain.LocalBranchName
	Empty
}

func (step *SetParent) Run(args RunArgs) error {
	err := args.Runner.Config.SetParent(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	return nil
}
