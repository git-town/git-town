package prompt

import (
	"fmt"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
)

// EnsureKnowsParentBranches asserts that the entire ancestry for all given branches
// is known to Git Town.
// Missing ancestry information is queried from the user.
func EnsureKnowsParentBranches(branchNames []string, repo *git.ProdRepo) error {
	for _, branchName := range branchNames {
		if repo.Config.IsMainBranch(branchName) || repo.Config.IsPerennialBranch(branchName) || repo.Config.HasParentBranch(branchName) {
			continue
		}
		err := AskForBranchAncestry(branchName, repo.Config.GetMainBranch(), repo)
		if err != nil {
			return err
		}
		if parentBranchHeaderShown {
			fmt.Println()
		}
	}
	return nil
}

// AskForBranchAncestry prompts the user for all unknown ancestors of the given branch.
func AskForBranchAncestry(branchName, defaultBranchName string, repo *git.ProdRepo) (err error) {
	current := branchName
	for {
		parent := repo.Config.GetParentBranch(current)
		if parent == "" { // nolint: nestif
			printParentBranchHeader(repo)
			parent, err = AskForBranchParent(current, defaultBranchName, repo)
			if err != nil {
				return err
			}
			if parent == perennialBranchOption {
				err = repo.Config.AddToPerennialBranches(current)
				if err != nil {
					return err
				}
				break
			}
			err = repo.Config.SetParentBranch(current, parent)
			if err != nil {
				return err
			}
		}
		if parent == repo.Config.GetMainBranch() || repo.Config.IsPerennialBranch(parent) {
			break
		}
		current = parent
	}
	return nil
}

// AskForBranchParent prompts the user for the parent of the given branch.
func AskForBranchParent(branchName, defaultBranchName string, repo *git.ProdRepo) (string, error) {
	choices, err := repo.Silent.LocalBranchesMainFirst()
	if err != nil {
		return "", err
	}
	filteredChoices := filterOutSelfAndDescendants(branchName, choices, repo)
	return askForBranch(askForBranchOptions{
		branchNames:       append([]string{perennialBranchOption}, filteredChoices...),
		prompt:            fmt.Sprintf(parentBranchPromptTemplate, branchName),
		defaultBranchName: defaultBranchName,
	}), nil
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

func filterOutSelfAndDescendants(branchName string, choices []string, repo *git.ProdRepo) (filteredChoices []string) {
	for _, choice := range choices {
		if choice == branchName || repo.Config.IsAncestorBranch(choice, branchName) {
			continue
		}
		filteredChoices = append(filteredChoices, choice)
	}
	return filteredChoices
}

func printParentBranchHeader(repo *git.ProdRepo) {
	if !parentBranchHeaderShown {
		parentBranchHeaderShown = true
		cli.Printf(parentBranchHeaderTemplate, repo.Config.GetMainBranch())
	}
}
