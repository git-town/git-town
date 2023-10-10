package step

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// DeleteLocalBranch deletes the branch with the given name.
type DeleteLocalBranch struct {
	Branch domain.LocalBranchName
	Force  bool
	Empty
}

func (step *DeleteLocalBranch) Run(args RunArgs) error {
	useForce := step.Force
	if !useForce {
		parent := args.Lineage.Parent(step.Branch)
		hasUnmergedCommits, err := args.Runner.Backend.BranchHasUnmergedCommits(step.Branch, parent.Location())
		if err != nil {
			return err
		}
		if hasUnmergedCommits {
			useForce = true
		}
	}
	return args.Runner.Frontend.DeleteLocalBranch(step.Branch, useForce)
}
