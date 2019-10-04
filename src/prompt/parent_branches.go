package prompt

import (
	"fmt"

	"github.com/Originate/git-town/src/cfmt"
	"github.com/Originate/git-town/src/git"
)

// EnsureKnowsParentBranches asserts that the entire ancestry for all given branches
// is known to Git Town.
// Missing ancestry information is queried from the user.
func EnsureKnowsParentBranches(branchNames []string) {
	for _, branchName := range branchNames {
		if git.IsMainBranch(branchName) || git.IsPerennialBranch(branchName) || git.HasParentBranch(branchName) {
			continue
		}
		AskForBranchAncestry(branchName, git.GetMainBranch())
		if parentBranchHeaderShown {
			fmt.Println()
		}
	}
}

// AskForBranchAncestry prompts the user for all unknown ancestors of the given branch
func AskForBranchAncestry(branchName, defaultBranchName string) {
	current := branchName
	for {
		parent := git.GetParentBranch(current)
		if parent == "" {
			printParentBranchHeader()
			parent = AskForBranchParent(current, defaultBranchName)
			if parent == perennialBranchOption {
				git.AddToPerennialBranches(current)
				break
			}
			git.SetParentBranch(current, parent)
		}
		if parent == git.GetMainBranch() || git.IsPerennialBranch(parent) {
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
var parentBranchPromptTemplate = "Please specify the parent branch of '%s':"
var perennialBranchOption = "<none> (perennial branch)"

func filterOutSelfAndDescendants(branchName string, choices []string) []string {
	result := []string{}
	for _, choice := range choices {
		if choice == branchName || git.IsAncestorBranch(choice, branchName) {
			continue
		}
		result = append(result, choice)
	}
	return result
}

func printParentBranchHeader() {
	if !parentBranchHeaderShown {
		parentBranchHeaderShown = true
		cfmt.Printf(parentBranchHeaderTemplate, git.GetMainBranch())
	}
}
