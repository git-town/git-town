package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
)

// DeleteLocalBranchStep deletes the branch with the given name.
type DeleteLocalBranchStep struct {
	Branch domain.LocalBranchName
	Parent domain.Location
	Force  bool
	EmptyStep
}

func (step *DeleteLocalBranchStep) Run(args RunArgs) error {
	hasUnmergedCommits, err := args.Runner.Backend.BranchHasUnmergedCommits(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	fmt.Println(hasUnmergedCommits)
	return args.Runner.Frontend.DeleteLocalBranch(step.Branch, step.Force || hasUnmergedCommits)
}
