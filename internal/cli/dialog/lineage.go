package dialog

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Lineage validates that the given lineage contains the ancestry for all given branches.
// Prompts missing lineage information from the user.
// Returns the new lineage and perennial branches to add to the config storage.
func Lineage(args LineageArgs) (LineageResult, dialogdomain.Exit, error) {
	additionalLineage := configdomain.NewLineage()
	additionalPerennials := gitdomain.LocalBranchNames{}
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
			if proposalSearcher, canSearchProposals := connector.(forgedomain.ProposalSearcher); canSearchProposals {
				proposals, _ := proposalSearcher.SearchProposals(branchToVerify)
				if len(proposals) == 1 {
					proposal := proposals[0]
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
			args.Config.NormalConfig.Lineage.Children(branchToVerify, args.Config.NormalConfig.Order)...,
		)
		noneEntry := SwitchBranchEntry{
			Branch:        messages.SetParentNoneOption,
			Indentation:   "",
			OtherWorktree: false,
			Type:          configdomain.BranchTypeFeatureBranch,
		}
		entriesArgs := NewSwitchBranchEntriesArgs{
			BranchInfos:       args.BranchInfos,
			BranchTypes:       []configdomain.BranchType{},
			BranchesAndTypes:  args.BranchesAndTypes,
			ExcludeBranches:   excludeBranches,
			Lineage:           args.Config.NormalConfig.Lineage,
			MainBranch:        Some(args.MainBranch),
			Order:             args.Config.NormalConfig.Order,
			Regexes:           []*regexp.Regexp{},
			ShowAllBranches:   true,
			UnknownBranchType: args.Config.NormalConfig.UnknownBranchType,
		}
		entriesAll := append(SwitchBranchEntries{noneEntry}, NewSwitchBranchEntries(entriesArgs)...)
		entriesArgs.ShowAllBranches = false
		entriesLocal := append(SwitchBranchEntries{noneEntry}, NewSwitchBranchEntries(entriesArgs)...)
		newParent, exit, err := SwitchBranch(SwitchBranchArgs{
			CurrentBranch:      None[gitdomain.LocalBranchName](),
			Cursor:             1, // select the "main branch" entry, below the "make perennial" entry
			DisplayBranchTypes: args.Config.NormalConfig.DisplayTypes,
			EntryData: EntryData{
				EntriesAll:      entriesAll,
				EntriesLocal:    entriesLocal,
				ShowAllBranches: false,
			},
			InputName:          fmt.Sprintf("parent-branch-for-%q", branchToVerify),
			Inputs:             args.Inputs,
			Title:              Some(fmt.Sprintf(messages.ParentBranchTitle, branchToVerify)),
			UncommittedChanges: false,
		})
		if err != nil || exit {
			if err != nil {
				err = fmt.Errorf(messages.NoTTYParentBranchMissing, branchToVerify, err)
			}
			return LineageResult{
				AdditionalLineage:    additionalLineage,
				AdditionalPerennials: additionalPerennials,
			}, exit, err
		}
		if newParent == messages.SetParentNoneOption {
			additionalPerennials = append(additionalPerennials, branchToVerify)
		} else {
			additionalLineage = additionalLineage.Set(branchToVerify, newParent)
			branchesToVerify = append(branchesToVerify, newParent)
		}
	}
	return LineageResult{
		AdditionalLineage:    additionalLineage,
		AdditionalPerennials: additionalPerennials,
	}, false, nil
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

type LineageResult struct {
	AdditionalLineage    configdomain.Lineage
	AdditionalPerennials gitdomain.LocalBranchNames
}
