package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/genericslice"
	"github.com/git-town/git-town/v9/src/git"
)

// EnterParent lets the user select a new parent for the given branch.
func EnterParent(branch, defaultParent string, lineage config.Lineage, branches git.BranchesSyncStatus) (string, error) {
	choices := genericslice.Hoist(branches.LocalBranches().BranchNames(), defaultParent)
	filteredChoices := filterOutSelfAndDescendants(branch, choices, lineage)
	return Select(SelectArgs{
		Options: append([]string{PerennialBranchOption}, filteredChoices...),
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
	PerennialBranchOption      = "<none> (perennial branch)"
)
