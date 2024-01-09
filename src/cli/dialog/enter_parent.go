package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
)

// EnterParent lets the user select a new parent for the given branch.
func EnterParent(branch, defaultParent gitdomain.LocalBranchName, lineage configdomain.Lineage, branches gitdomain.BranchInfos) (gitdomain.LocalBranchName, error) {
	choices := branches.LocalBranches().Names()
	choices = slice.Hoist(choices, defaultParent)
	filteredChoices := filterOutSelfAndDescendants(branch, choices, lineage)
	choice, err := Select(SelectArgs{
		Options: append([]string{PerennialBranchOption}, filteredChoices.Strings()...),
		Message: fmt.Sprintf(parentBranchPromptTemplate, branch),
		Default: defaultParent.String(),
	})
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), err
	}
	return gitdomain.NewLocalBranchName(choice), nil
}

func filterOutSelfAndDescendants(branch gitdomain.LocalBranchName, choices gitdomain.LocalBranchNames, lineage configdomain.Lineage) gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames{}
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
