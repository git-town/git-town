package step

import "github.com/git-town/git-town/v9/src/domain"

type RemoveBranchFromLineage struct {
	Branch domain.LocalBranchName
	Empty
}

func (step *RemoveBranchFromLineage) Run(args RunArgs) error {
	// Ignoring errors removing the config here since the config entry might not exist,
	// for example when pruning perennial branches or branches with unknown ancestry.
	parent := args.Lineage.Parent(step.Branch)
	for _, child := range args.Lineage.Children(step.Branch) {
		if parent.IsEmpty() {
			_ = args.Runner.Backend.Config.RemoveParent(child)
		} else {
			_ = args.Runner.Backend.Config.SetParent(child, parent)
		}
	}
	_ = args.Runner.Backend.Config.RemoveParent(step.Branch)
	args.RemoveBranchFromLineage(step.Branch)
	return nil
}
