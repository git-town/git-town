package step

import "github.com/git-town/git-town/v9/src/domain"

type RemoveBranchFromLineage struct {
	Branch domain.LocalBranchName
	Empty
}

func (step *RemoveBranchFromLineage) Run(args RunArgs) error {
	args.RemoveBranchFromLineage(step.Branch)
	return nil
}
