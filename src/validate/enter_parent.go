package validate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/stringslice"
)

// EnterParent lets the user select a new parent for the given branch.
func EnterParent(branch, defaultParent string, lineage config.Lineage, branches git.BranchesSyncStatus) (string, error) {
	choices := stringslice.Hoist(branches.LocalBranches().BranchNames(), defaultParent)
	filteredChoices := filterOutSelfAndDescendants(branch, choices, lineage)
	return dialog.Select(dialog.SelectArgs{
		Options: append([]string{perennialBranchOption}, filteredChoices...),
		Message: fmt.Sprintf(parentBranchPromptTemplate, branch),
		Default: defaultParent,
	})
}

func filterOutSelfAndDescendants(branch string, choices []string, lineage config.Lineage) []string {
	result := []string{}
	for _, choice := range choices {
		if choice == branch || lineage.IsAncestor(branch, choice) {
			continue
		}
		result = append(result, choice)
	}
	return result
}

const (
	parentBranchPromptTemplate = "Please specify the parent branch of %q:"
	perennialBranchOption      = "<none> (perennial branch)"
)
