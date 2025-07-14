package dialog

import (
	"fmt"
	"slices"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	"github.com/git-town/git-town/v21/internal/messages"
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

The main branch is automatically perennial
and therefore always selected in this list.
`
)

// PerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func PerennialBranches(localBranches gitdomain.LocalBranchNames, oldPerennialBranches gitdomain.LocalBranchNames, mainBranch gitdomain.LocalBranchName, inputs dialogcomponents.TestInput) (gitdomain.LocalBranchNames, dialogdomain.Exit, error) {
	perennialCandidates := localBranches.AppendAllMissing(oldPerennialBranches...)
	if len(perennialCandidates) < 2 {
		return gitdomain.LocalBranchNames{}, false, nil
	}
	entries := make(list.Entries[gitdomain.LocalBranchName], len(perennialCandidates))
	for b, branch := range perennialCandidates {
		isMain := branch == mainBranch
		entries[b] = list.Entry[gitdomain.LocalBranchName]{
			Data:     branch,
			Disabled: isMain,
			Text:     branch.String(),
		}
	}
	selections := slice.FindMany(perennialCandidates, oldPerennialBranches)
	selections = append(selections, slices.Index(perennialCandidates, mainBranch))
	selectedBranchesList, exit, err := dialogcomponents.CheckList(entries, selections, perennialBranchesTitle, PerennialBranchesHelp, inputs)
	selectedBranches := gitdomain.LocalBranchNames(selectedBranchesList)
	selectedBranches = selectedBranches.Remove(mainBranch)
	selectionText := selectedBranches.Join(", ")
	if selectionText == "" {
		selectionText = "(none)"
	}
	fmt.Printf(messages.PerennialBranches, dialogcomponents.FormattedSelection(selectionText, exit))
	return selectedBranches, exit, err
}
