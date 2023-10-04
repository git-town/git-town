package step

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// DeleteLocalBranch deletes the branch with the given name.
type DeleteLocalBranch struct {
	Branch domain.LocalBranchName
	Parent domain.Location
	Force  bool
	Empty
}

func (step *DeleteLocalBranch) Run(args RunArgs) error {
	hasUnmergedCommits, err := args.Runner.Backend.BranchHasUnmergedCommits(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	return args.Runner.Frontend.DeleteLocalBranch(step.Branch, step.Force || hasUnmergedCommits)
}
