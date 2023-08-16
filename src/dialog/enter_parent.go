package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/slice"
)

// EnterParent lets the user select a new parent for the given branch.
func EnterParent(branch, defaultParent domain.LocalBranchName, lineage config.Lineage, branches git.BranchesSyncStatus) (domain.LocalBranchName, error) {
	choices := slice.Hoist(branches.LocalBranches().LocalBranchNames(), defaultParent)
	filteredChoices := filterOutSelfAndDescendants(branch, choices, lineage)
	choice, err := Select(SelectArgs{
		Options: append([]string{PerennialBranchOption}, filteredChoices.Strings()...),
		Message: fmt.Sprintf(parentBranchPromptTemplate, branch),
		Default: defaultParent.String(),
	})
	if err != nil {
		return domain.LocalBranchName{}, err
	}
	return domain.NewLocalBranchName(choice), nil
}

func filterOutSelfAndDescendants(branch domain.LocalBranchName, choices domain.LocalBranchNames, lineage config.Lineage) domain.LocalBranchNames {
	result := domain.LocalBranchNames{}
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
