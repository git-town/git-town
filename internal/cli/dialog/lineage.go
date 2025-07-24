package dialog

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Lineage validates that the given lineage contains the ancestry for all given branches.
// Prompts missing lineage information from the user.
// Returns the new lineage and perennial branches to add to the config storage.
func Lineage(args LineageArgs) (additionalLineage configdomain.Lineage, additionalPerennials gitdomain.LocalBranchNames, exit dialogdomain.Exit, err error) {
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
		if slices.Contains(args.Config.NormalConfig.PerennialBranches, branchToVerify) {
			continue
		}
		if parent, hasParent := args.Config.NormalConfig.Lineage.Parent(branchToVerify).Get(); hasParent {
			branchesToVerify = append(branchesToVerify, parent)
			continue
		}
		// look for parent in proposals
		if connector, hasConnector := args.Connector.Get(); hasConnector {
			if searchProposals, canSearchProposals := connector.SearchProposalFn().Get(); canSearchProposals {
				proposalOpt, _ := searchProposals(branchToVerify)
				if proposal, hasProposal := proposalOpt.Get(); hasProposal {
					parent := proposal.Data.Data().Target
					additionalLineage = additionalLineage.Set(branchToVerify, parent)
					branchesToVerify = append(branchesToVerify, parent)
					continue
				}
			}
		}
		// ask for parent
		excludeBranches := append(
			gitdomain.LocalBranchNames{branchToVerify},
			args.Config.NormalConfig.Lineage.Children(branchToVerify)...,
		)
		entries := NewSwitchBranchEntries(NewSwitchBranchEntriesArgs{
			BranchInfos:       args.BranchInfos,
			BranchTypes:       []configdomain.BranchType{},
			BranchesAndTypes:  args.BranchesAndTypes,
			ExcludeBranches:   excludeBranches,
			Lineage:           args.Config.NormalConfig.Lineage,
			Regexes:           []*regexp.Regexp{},
			ShowAllBranches:   true,
			UnknownBranchType: args.Config.NormalConfig.UnknownBranchType,
		})
		noneEntry := SwitchBranchEntry{
			Branch:        messages.SetParentNoneOption,
			Indentation:   "",
			OtherWorktree: false,
			Type:          configdomain.BranchTypeFeatureBranch,
		}
		entries = append(SwitchBranchEntries{noneEntry}, entries...)
		selectedBranch, exit, err := SwitchBranch(SwitchBranchArgs{
			CurrentBranch:      None[gitdomain.LocalBranchName](),
			Cursor:             1,
			DisplayBranchTypes: false,
			Entries:            entries,
			InputName:          fmt.Sprintf("parent-branch-for-%q", branchToVerify),
			Inputs:             args.Inputs,
			Title:              Some(fmt.Sprintf(ParentBranchTitleTemplate, branchToVerify)),
			UncommittedChanges: false,
		})
		if err != nil || exit {
			return additionalLineage, additionalPerennials, exit, err
		}
		var outcome ParentOutcome
		if selectedBranch == messages.SetParentNoneOption {
			outcome = ParentOutcomePerennialBranch
		} else {
			outcome = ParentOutcomeSelectedParent
		}
		switch outcome {
		case ParentOutcomeExit:
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
	BranchInfos      gitdomain.BranchInfos
	BranchesAndTypes configdomain.BranchesAndTypes
	BranchesToVerify gitdomain.LocalBranchNames
	Config           config.UnvalidatedConfig
	Connector        Option[forgedomain.Connector]
	DefaultChoice    gitdomain.LocalBranchName
	Inputs           dialogcomponents.Inputs
	LocalBranches    gitdomain.LocalBranchNames
	MainBranch       gitdomain.LocalBranchName
}
