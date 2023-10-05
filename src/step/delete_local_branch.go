package step

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// DeleteLocalBranch deletes the branch with the given name.
type DeleteLocalBranch struct {
	Branch domain.LocalBranchName
	Parent domain.Location
	Force  bool // TODO: either always determine at the call site whether to force-delete or not and strictly follow this flag and remove the call to BranchHasUnmergedCommits, or remove this flag and keep calling BranchHasUnmergedCommits - but don't keep doing both.
	Empty
}

func (step *DeleteLocalBranch) Run(args RunArgs) error {
	useForce := step.Force
	if !useForce {
		hasUnmergedCommits, err := args.Runner.Backend.BranchHasUnmergedCommits(step.Branch, step.Parent)
		if err != nil {
			return err
		}
		if hasUnmergedCommits {
			useForce = true
		}
	}
	return args.Runner.Frontend.DeleteLocalBranch(step.Branch, useForce)
}
