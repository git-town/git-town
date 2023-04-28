package validate

import (
	"fmt"

	"github.com/git-town/git-town/v8/src/dialog"
	"github.com/git-town/git-town/v8/src/git"
)

// EnterParent lets the user select a new parent for the given branch.
func EnterParent(branch, defaultParent string, backend *git.BackendCommands) (string, error) {
	choices, err := backend.LocalBranchesMainFirst(defaultParent)
	if err != nil {
		return "", err
	}
	filteredChoices := filterOutSelfAndDescendants(branch, choices, backend.RepoConfig)
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
