package validate

import (
	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

// Lineage validates that the given lineage contains the ancestry for all given branches.
// Prompts missing lineage information from the user and updates persisted lineage as needed.
// Returns the validated Lineage.
func Lineage(args LineageArgs) (additionalLineage configdomain.Lineage, additionalPerennials gitdomain.LocalBranchNames, aborted bool, err error) {
	// step 1: determine all branches for which the parent must be known
	branchesToVerify := args.BranchesToVerify

	// step 2: for each branch: check the ancestor
	for _, branchToVerify := range args.BranchesToVerify {
		parent, hasParent := args.Config.Lineage.Parent(branchToVerify).Get()
		if hasParent {
			branchesToVerify = append(branchesToVerify, parent)
			continue
		}
		outcome, selectedBranch, err := dialog.Parent(dialog.ParentArgs{
			Branch:          branchToVerify,
			DefaultChoice:   args.DefaultChoice,
			DialogTestInput: args.DialogTestInputs.Next(),
			Lineage:         args.Config.Lineage,
			LocalBranches:   args.LocalBranches,
			MainBranch:      args.MainBranch,
		})
		if err != nil {
			return additionalLineage, additionalPerennials, false, err
		}
		switch outcome {
		case dialog.ParentOutcomeAborted:
			return additionalLineage, additionalPerennials, true, nil
		case dialog.ParentOutcomePerennialBranch:
			additionalPerennials = append(additionalPerennials, branchToVerify)
		case dialog.ParentOutcomeSelectedParent:
			additionalLineage[branchToVerify] = selectedBranch
		}
		if args.Config.IsMainOrPerennialBranch(selectedBranch) {
			break
		}
	}
	return additionalLineage, additionalPerennials, false, nil
}

type LineageArgs struct {
	BranchesToVerify gitdomain.LocalBranchNames
	Config           *configdomain.UnvalidatedConfig
	DefaultChoice    gitdomain.LocalBranchName
	DialogTestInputs *components.TestInputs
	LocalBranches    gitdomain.LocalBranchNames
	MainBranch       gitdomain.LocalBranchName
}
