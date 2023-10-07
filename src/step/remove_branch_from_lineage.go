package step

import "github.com/git-town/git-town/v9/src/domain"

type RemoveBranchFromLineage struct {
	Branch domain.LocalBranchName
	Empty
}

func (step *RemoveBranchFromLineage) Run(args RunArgs) error {
	parent := args.Lineage.Parent(step.Branch)
	for _, child := range args.Lineage.Children(step.Branch) {
		if parent.IsEmpty() {
			args.Runner.Backend.Config.RemoveParent(child)
		} else {
			err := args.Runner.Backend.Config.SetParent(child, parent)
			if err != nil {
				return err
			}
		}
	}
	// Ignoring errors removing the config here since the config entry might not exist,
	// for example when removing perennial branches or branches with unknown ancestry.
	args.Runner.Backend.Config.RemoveParent(step.Branch)
	args.RemoveBranchFromLineage(step.Branch)
	return nil
}
