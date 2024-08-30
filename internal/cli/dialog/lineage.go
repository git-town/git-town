package dialog

import (
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
)

// Lineage validates that the given lineage contains the ancestry for all given branches.
// Prompts missing lineage information from the user.
// Returns the new lineage and perennial branches to add to the config storage.
func Lineage(args LineageArgs) (additionalLineage configdomain.Lineage, additionalPerennials gitdomain.LocalBranchNames, aborted bool, err error) {
	additionalLineage = configdomain.NewLineage()
	branchesToVerify := args.BranchesToVerify
	for i := 0; i < len(branchesToVerify); i++ {
		branchToVerify := branchesToVerify[i]
		if !args.Config.MustKnowParent(branchToVerify) {
			continue
		}
		if parent, hasParent := args.Config.Lineage.Parent(branchToVerify).Get(); hasParent {
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
			additionalLineage.Add(branchToVerify, selectedBranch)
			branchesToVerify = append(branchesToVerify, selectedBranch)
		}
	}
	return additionalLineage, additionalPerennials, false, nil
}

type LineageArgs struct {
	BranchesToVerify gitdomain.LocalBranchNames
	Config           configdomain.UnvalidatedConfig
	DefaultChoice    gitdomain.LocalBranchName
	DialogTestInputs components.TestInputs
	LocalBranches    gitdomain.LocalBranchNames
	MainBranch       gitdomain.LocalBranchName
}
