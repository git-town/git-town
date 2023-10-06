package step

import "github.com/git-town/git-town/v9/src/domain"

type RemoveBranchFromLineage struct {
	Branch domain.LocalBranchName
	Empty
}

func (step *RemoveBranchFromLineage) Run(args RunArgs) error {
	err := args.Runner.Backend.Config.RemoveParent(step.Branch)
	args.RemoveBranchFromLineage(step.Branch)
	return err
}
