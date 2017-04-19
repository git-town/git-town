package prompt

import (
	"errors"
	"fmt"

	"github.com/Originate/git-town/lib/git"
	"github.com/fatih/color"
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
var parentBranchPromptTemplate = "Please specify the parent branch of %s by name or number (default: %s): "

func askForBranchAncestry(branchName string) {
	current := branchName
	for {
		parent := git.GetParentBranch(current)
		if parent == "" {
			printParentBranchHeader()
			parent = askForBranch(branchPromptConfig{
				branchNames: git.GetLocalBranchesWithMainBranchFirst(),
				prompt:      getParentBranchPrompt(current),
				validate: func(branchName string) error {
					return validateParentBranch(current, branchName)
				},
			})
			if parent == "" {
				parent = git.GetMainBranch()
			}
			git.SetParentBranch(current, parent)
		}
		if parent == git.GetMainBranch() || git.IsPerennialBranch(parent) {
			break
		}
		current = parent
	}
}

func validateParentBranch(branchName string, parent string) error {
	if branchName == parent {
		return errors.New(fmt.Sprintf("'%s' cannot be the parent of itself", parent))
	}
	if branchName != "" && git.IsAncestorBranch(parent, branchName) {
		return errors.New(fmt.Sprintf("Nested branch loop detected: '%s' is an ancestor of '%s'", branchName, parent))
	}
	return nil
}

func printParentBranchHeader() {
	if !parentBranchHeaderShown {
		parentBranchHeaderShown = true
		fmt.Printf(parentBranchHeaderTemplate, git.GetMainBranch())
		printNumberedBranches(git.GetLocalBranchesWithMainBranchFirst())
		fmt.Println()
	}
}

func getParentBranchPrompt(branchName string) string {
	coloredBranchName := color.New(color.Bold).Add(color.FgCyan).Sprintf(branchName)
	return fmt.Sprintf(parentBranchPromptTemplate, coloredBranchName, git.GetMainBranch())
}
