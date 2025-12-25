package dialog

import (
	"fmt"
	"slices"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/slice"
	"github.com/git-town/git-town/v22/internal/messages"
)

const (
	perennialBranchesTitle = `Perennial branches`
	PerennialBranchesHelp  = `
Perennial branches are long-lived branches
that aren't shipped and don't have parent branches.
They typically represent environments like
development, staging, qa, or production.

For more flexible configuration,
you can also use the "perennial-regex" setting
to match branch names dynamically.

Always selected in this list are:
- the main branch because it is automatically perennial
- perennial branches defined in the config file
- perennial branches from the global Git configuration

`
)

// PerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func PerennialBranches(args PerennialBranchesArgs) (gitdomain.LocalBranchNames, dialogdomain.Exit, error) {
	perennialCandidates := args.LocalBranches.AppendAllMissing(args.ImmutableGitPerennials).AppendAllMissing(args.LocalGitPerennials)
	if len(perennialCandidates) < 2 {
		// there is always the main branch in this list, so if that's the only one there is no branch to select --> don't display the dialog
		return gitdomain.LocalBranchNames{}, false, nil
	}
	entries := make(list.Entries[gitdomain.LocalBranchName], len(perennialCandidates))
	for b, branch := range perennialCandidates {
		isImmutablePerennial := args.ImmutableGitPerennials.Contains(branch)
		entries[b] = list.Entry[gitdomain.LocalBranchName]{
			Data:     branch,
			Disabled: isImmutablePerennial,
			Text:     branch.String(),
		}
	}
	selections := []int{slices.Index(perennialCandidates, args.MainBranch)}
	selections = append(selections, slice.FindMany(perennialCandidates, args.ImmutableGitPerennials)...)
	selections = append(selections, slice.FindMany(perennialCandidates, args.LocalGitPerennials)...)
	selectedBranchesList, exit, err := dialogcomponents.CheckList(entries, selections, perennialBranchesTitle, PerennialBranchesHelp, args.Inputs, "perennial-branches")
	selectedBranches := gitdomain.LocalBranchNames(selectedBranchesList)
	selectedBranches = selectedBranches.Remove(args.ImmutableGitPerennials...)
	selectionText := selectedBranches.Join(", ")
	if selectionText == "" {
		selectionText = messages.DialogResultNone
	}
	fmt.Printf(messages.PerennialBranches, dialogcomponents.FormattedSelection(selectionText, exit))
	return selectedBranches, exit, err
}

type PerennialBranchesArgs struct {
	ImmutableGitPerennials gitdomain.LocalBranchNames // perennial branches defined in the config file and the global Git metadata
	Inputs                 dialogcomponents.Inputs
	LocalBranches          gitdomain.LocalBranchNames
	LocalGitPerennials     gitdomain.LocalBranchNames
	MainBranch             gitdomain.LocalBranchName
}
