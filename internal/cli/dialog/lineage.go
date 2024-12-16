package dialog

import (
	"slices"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/hosting/hostingdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

// Lineage validates that the given lineage contains the ancestry for all given branches.
// Prompts missing lineage information from the user.
// Returns the new lineage and perennial branches to add to the config storage.
func Lineage(args LineageArgs) (additionalLineage configdomain.Lineage, additionalPerennials gitdomain.LocalBranchNames, aborted bool, err error) {
	additionalLineage = configdomain.NewLineage()
	branchesToVerify := args.BranchesToVerify
	for i := 0; i < len(branchesToVerify); i++ {
		branchToVerify := branchesToVerify[i]
		branchType, hasBranchType := args.BranchesAndTypes[branchToVerify]
		if hasBranchType && !branchType.MustKnowParent() {
			continue
		}
		// If the main branch isn't local, it isn't in args.BranchesAndTypes.
		// We therefore exclude it manually here.
		if branchToVerify == args.MainBranch {
			continue
		}
		// If a perennial branch isn't local, it isn't in args.BranchesAndTypes.
		// We therefore exclude them manually here.
		if slices.Contains(args.PerennialBranches, branchToVerify) {
			continue
		}
		if parent, hasParent := args.Lineage.Parent(branchToVerify).Get(); hasParent {
			branchesToVerify = append(branchesToVerify, parent)
			continue
		}
		// look for parent in proposals
		if connector, hasConnector := args.Connector.Get(); hasConnector {
			if searchProposals, canSearchProposals := connector.SearchProposalFn().Get(); canSearchProposals {
				proposalOpt, _ := searchProposals(branchToVerify)
				if proposal, hasProposal := proposalOpt.Get(); hasProposal {
					parent := proposal.Target
					additionalLineage = additionalLineage.Set(branchToVerify, parent)
					branchesToVerify = append(branchesToVerify, parent)
					continue
				}
			}
		}
		// ask for parent
		outcome, selectedBranch, err := Parent(ParentArgs{
			Branch:          branchToVerify,
			DefaultChoice:   args.DefaultChoice,
			DialogTestInput: args.DialogTestInputs.Next(),
			Lineage:         args.Lineage,
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
			additionalLineage = additionalLineage.Set(branchToVerify, selectedBranch)
			branchesToVerify = append(branchesToVerify, selectedBranch)
		}
	}
	return additionalLineage, additionalPerennials, false, nil
}

type LineageArgs struct {
	BranchesAndTypes  configdomain.BranchesAndTypes
	BranchesToVerify  gitdomain.LocalBranchNames
	Connector         Option[hostingdomain.Connector]
	DefaultChoice     gitdomain.LocalBranchName
	DialogTestInputs  components.TestInputs
	Lineage           configdomain.Lineage
	LocalBranches     gitdomain.LocalBranchNames
	MainBranch        gitdomain.LocalBranchName
	PerennialBranches gitdomain.LocalBranchNames
}
