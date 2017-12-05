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
		if git.IsMainBranch(branchName) || git.IsPerennialBranch(branchName) || git.HasCompiledAncestorBranches(branchName) {
			continue
		}
		askForBranchAncestry(branchName)
		ancestors := git.CompileAncestorBranches(branchName)
		git.SetAncestorBranches(branchName, ancestors)

		if parentBranchHeaderShown {
			fmt.Println()
		}
	}
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

func askForBranchAncestry(branchName string) {
	current := branchName
	choices := git.GetLocalBranchesWithMainBranchFirst()
	for {
		parent := git.GetParentBranch(current)
		if parent == "" {
			printParentBranchHeader()
			filteredChoices := filterOutSelfAndDescendants(current, choices)
			parent = askForBranch(askForBranchOptions{
				branchNames:       append([]string{perennialBranchOption}, filteredChoices...),
				prompt:            fmt.Sprintf(parentBranchPromptTemplate, current),
				defaultBranchName: git.GetMainBranch(),
			})
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
