package dialog

import (
	"fmt"

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

`
)

// PerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func PerennialBranches(localBranches gitdomain.LocalBranchNames, oldPerennialBranches gitdomain.LocalBranchNames, mainBranch gitdomain.LocalBranchName, inputs dialogcomponents.TestInput) (gitdomain.LocalBranchNames, dialogdomain.Exit, error) {
	perennialCandidates := localBranches.Remove(mainBranch).AppendAllMissing(oldPerennialBranches...)
	if len(perennialCandidates) == 0 {
		return gitdomain.LocalBranchNames{}, false, nil
	}
	entries := list.NewEntries(perennialCandidates...)
	selections := slice.FindMany(perennialCandidates, oldPerennialBranches)
	selectedBranchesList, exit, err := dialogcomponents.CheckList(entries, selections, perennialBranchesTitle, PerennialBranchesHelp, inputs)
	selectedBranches := gitdomain.LocalBranchNames(selectedBranchesList)
	selectionText := selectedBranches.Join(", ")
	if selectionText == "" {
		selectionText = "(none)"
	}
	fmt.Printf(messages.PerennialBranches, dialogcomponents.FormattedSelection(selectionText, exit))
	return selectedBranches, exit, err
}
