package dialog

import (
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

// Lineage validates that the given lineage contains the ancestry for all given branches.
// Prompts missing lineage information from the user and updates persisted lineage as needed.
// Returns the validated Lineage.
func Lineage(args LineageArgs) (additionalLineage configdomain.Lineage, additionalPerennials gitdomain.LocalBranchNames, aborted bool, err error) {
	additionalLineage = make(configdomain.Lineage)
	branchesToVerify := args.BranchesToVerify
	for _, branchToVerify := range args.BranchesToVerify {
		parent, hasParent := args.Config.Lineage.Parent(branchToVerify).Get()
		if hasParent {
			branchesToVerify = append(branchesToVerify, parent)
			continue
		}
		outcome, selectedBranch, err := Parent(ParentArgs{
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
		case ParentOutcomeAborted:
			return additionalLineage, additionalPerennials, true, nil
		case ParentOutcomePerennialBranch:
			additionalPerennials = append(additionalPerennials, branchToVerify)
		case ParentOutcomeSelectedParent:
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
	Config           configdomain.UnvalidatedConfig
	DefaultChoice    gitdomain.LocalBranchName
	DialogTestInputs *components.TestInputs
	LocalBranches    gitdomain.LocalBranchNames
	MainBranch       gitdomain.LocalBranchName
}
