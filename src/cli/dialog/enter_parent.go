package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
)

// EnterParent lets the user select a new parent for the given branch.
func EnterParent(branch, defaultParent domain.LocalBranchName, lineage configdomain.Lineage, branches domain.BranchInfos) (domain.LocalBranchName, error) {
	choices := branches.LocalBranches().Names()
	slice.Hoist(&choices, defaultParent)
	filteredChoices := filterOutSelfAndDescendants(branch, choices, lineage)
	choice, err := Select(SelectArgs{
		Options: append([]string{PerennialBranchOption}, filteredChoices.Strings()...),
		Message: fmt.Sprintf(parentBranchPromptTemplate, branch),
		Default: defaultParent.String(),
	})
	if err != nil {
		return domain.EmptyLocalBranchName(), err
	}
	return domain.NewLocalBranchName(choice), nil
}

func filterOutSelfAndDescendants(branch domain.LocalBranchName, choices domain.LocalBranchNames, lineage configdomain.Lineage) domain.LocalBranchNames {
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
