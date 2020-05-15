package prompt

import (
	"fmt"

	"github.com/git-town/git-town/src/cfmt"
	"github.com/git-town/git-town/src/git"
)

// EnsureKnowsParentBranches asserts that the entire ancestry for all given branches
// is known to Git Town.
// Missing ancestry information is queried from the user.
func EnsureKnowsParentBranches(branchNames []string) {
	for _, branchName := range branchNames {
		if git.Config().IsMainBranch(branchName) || git.Config().IsPerennialBranch(branchName) || git.Config().HasParentBranch(branchName) {
			continue
		}
		AskForBranchAncestry(branchName, git.Config().GetMainBranch())
		if parentBranchHeaderShown {
			fmt.Println()
		}
	}
}

// AskForBranchAncestry prompts the user for all unknown ancestors of the given branch
func AskForBranchAncestry(branchName, defaultBranchName string) {
	current := branchName
	for {
		parent := git.Config().GetParentBranch(current)
		if parent == "" {
			printParentBranchHeader()
			parent = AskForBranchParent(current, defaultBranchName)
			if parent == perennialBranchOption {
				git.Config().AddToPerennialBranches(current)
				break
			}
			git.Config().SetParentBranch(current, parent)
		}
		if parent == git.Config().GetMainBranch() || git.Config().IsPerennialBranch(parent) {
			break
		}
		current = parent
	}
}

// AskForBranchParent prompts the user for the parent of the given branch
func AskForBranchParent(branchName, defaultBranchName string) string {
	choices := git.GetLocalBranchesWithMainBranchFirst()
	filteredChoices := filterOutSelfAndDescendants(branchName, choices)
	return askForBranch(askForBranchOptions{
		branchNames:       append([]string{perennialBranchOption}, filteredChoices...),
		prompt:            fmt.Sprintf(parentBranchPromptTemplate, branchName),
		defaultBranchName: defaultBranchName,
	})
}

// Helpers

var parentBranchHeaderShown = false
var parentBranchHeaderTemplate = `
Feature branches can be branched directly off
%s or from other feature branches.

The former allows to develop and ship features completely independent of each other.
The latter allows to build on top of currently unshipped features.

`
var parentBranchPromptTemplate = "Please specify the parent branch of %q:"
var perennialBranchOption = "<none> (perennial branch)"

func filterOutSelfAndDescendants(branchName string, choices []string) []string {
	result := []string{}
	for _, choice := range choices {
		if choice == branchName || git.Config().IsAncestorBranch(choice, branchName) {
			continue
		}
		result = append(result, choice)
	}
	return result
}

func printParentBranchHeader() {
	if !parentBranchHeaderShown {
		parentBranchHeaderShown = true
		cfmt.Printf(parentBranchHeaderTemplate, git.Config().GetMainBranch())
	}
}
