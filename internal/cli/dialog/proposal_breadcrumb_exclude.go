package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v23/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v23/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v23/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/messages"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

const (
	proposalBreadcrumbExcludeBranchesTitle = `Proposal Breadcrumb Exclude Branches`
	proposalBreadcrumbExcludeBranchesHelp  = `
	Which branch types should proposal breadcrumbs exclude?

	See https://www.git-town.com/preferences/breadcrumb-exclude-branches.html for details.
	`
)

func ProposalBreadcrumbExclude(args Args[configdomain.ProposalBreadcrumbExclude]) (Option[configdomain.ProposalBreadcrumbExclude], dialogdomain.Exit, error) {
	entries := list.NewEntries(configdomain.AllBranchTypes()...)
	selection, exit, err := dialogcomponents.CheckList(
		entries,
		proposalBreadcrumbExcludeBranchDefaultSelections(entries, args.Local.Or(args.Global)),
		proposalBreadcrumbExcludeBranchesTitle,
		proposalBreadcrumbExcludeBranchesHelp,
		args.Inputs,
		args.Interactive,
		"proposal-breadcrumb-exclude-branches",
	)
	result := Some(configdomain.NewProposalBreadcrumbExclude(selection...))
	if args.Global.Equal(result) {
		result = None[configdomain.ProposalBreadcrumbExclude]()
	}

	fmt.Printf(messages.ProposalBreadcrumbExclude, formatProposalBreadcrumbExclude(result, args.Global.IsSome(), exit))
	return result, exit, err
}

func proposalBreadcrumbExcludeBranchDefaultSelections(entries list.Entries[configdomain.BranchType], value Option[configdomain.ProposalBreadcrumbExclude]) []int {
	excludedBranches, hasExcludedBranches := value.Get()
	if !hasExcludedBranches {
		return []int{}
	}
	result := []int{}
	for entryNumber, entry := range entries {
		if excludedBranches.Contains(entry.Data) {
			result = append(result, entryNumber)
		}
	}
	return result
}

func formatProposalBreadcrumbExclude(selection Option[configdomain.ProposalBreadcrumbExclude], hasGlobal bool, exit dialogdomain.Exit) string {
	if value, hasValue := selection.Get(); hasValue {
		return dialogcomponents.FormattedOption(Some(value), hasGlobal, exit)
	}
	return dialogcomponents.FormattedOption(None[configdomain.ProposalBreadcrumbExclude](), hasGlobal, exit)
}
