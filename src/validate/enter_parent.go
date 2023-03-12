package validate

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
)

// EnterParent prompts the user for the parent of the given branch.
func EnterParent(branch, defaultBranch string, repo *git.ProdRepo) (string, error) {
	choices, err := repo.Silent.LocalBranchesMainFirst()
	if err != nil {
		return "", err
	}
	filteredChoices := filterOutSelfAndDescendants(branch, choices, repo)
	return dialog.Select(dialog.SelectArgs{
		Options: append([]string{perennialBranchOption}, filteredChoices...),
		Message: fmt.Sprintf(parentBranchPromptTemplate, branch),
		Default: defaultBranch,
	})
}

func filterOutSelfAndDescendants(branch string, choices []string, repo *git.ProdRepo) []string {
	result := []string{}
	for _, choice := range choices {
		if choice == branch || repo.Config.IsAncestorBranch(choice, branch) {
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
