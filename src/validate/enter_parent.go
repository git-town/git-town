package validate

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
)

// EnterParent lets the user select a new parent for the given branch.
func EnterParent(branch, defaultParent string, repo *git.InternalCommands) (string, error) {
	choices, err := repo.LocalBranchesMainFirst(defaultParent)
	if err != nil {
		return "", err
	}
	filteredChoices := filterOutSelfAndDescendants(branch, choices, repo.Config)
	return dialog.Select(dialog.SelectArgs{
		Options: append([]string{perennialBranchOption}, filteredChoices...),
		Message: fmt.Sprintf(parentBranchPromptTemplate, branch),
		Default: defaultParent,
	})
}

func filterOutSelfAndDescendants(branch string, choices []string, config *git.RepoConfig) []string {
	result := []string{}
	for _, choice := range choices {
		if choice == branch || config.IsAncestorBranch(choice, branch) {
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
