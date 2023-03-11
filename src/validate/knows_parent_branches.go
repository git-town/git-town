package validate

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
)

// KnowsBranchesAncestry asserts that the entire ancestry for all given branches
// is known to Git Town.
// Missing ancestry information is queried from the user.
func KnowsBranchesAncestry(branches []string, repo *git.ProdRepo) error {
	for _, branch := range branches {
		if repo.Config.IsMainBranch(branch) || repo.Config.IsPerennialBranch(branch) || repo.Config.HasParentBranch(branch) {
			continue
		}
		headerShown, err := KnowsBranchAncestry(branch, repo.Config.MainBranch(), repo)
		if err != nil {
			return err
		}
		if headerShown {
			fmt.Println()
		}
	}
	return nil
}

// KnowsBranchAncestry prompts the user for all unknown ancestors of the given branch.
func KnowsBranchAncestry(branch, defaultBranch string, repo *git.ProdRepo) (headerShown bool, err error) { //nolint:nonamedreturns // return value names are useful here
	currentBranch := branch
	for {
		parent := repo.Config.ParentBranch(currentBranch)
		if parent == "" { //nolint:nestif
			if !headerShown {
				printParentBranchHeader(repo)
				headerShown = true
			}
			parent, err = AskForParent(currentBranch, defaultBranch, repo)
			if err != nil {
				return
			}
			if parent == perennialBranchOption {
				err = repo.Config.AddToPerennialBranches(currentBranch)
				if err != nil {
					return
				}
				break
			}
			err = repo.Config.SetParent(currentBranch, parent)
			if err != nil {
				return
			}
		}
		if parent == repo.Config.MainBranch() || repo.Config.IsPerennialBranch(parent) {
			break
		}
		currentBranch = parent
	}
	return
}

// AskForParent prompts the user for the parent of the given branch.
func AskForParent(branch, defaultBranch string, repo *git.ProdRepo) (string, error) {
	choices, err := repo.Silent.LocalBranchesMainFirst()
	if err != nil {
		return "", err
	}
	filteredChoices := filterOutSelfAndDescendants(branch, choices, repo)
	return dialog.AskForBranch(dialog.AskForBranchOptions{
		Branches:      append([]string{perennialBranchOption}, filteredChoices...),
		Prompt:        fmt.Sprintf(parentBranchPromptTemplate, branch),
		DefaultBranch: defaultBranch,
	})
}

// Helpers

const parentBranchHeaderTemplate string = `
Feature branches can be branched directly off
%s or from other feature branches.

The former allows to develop and ship features completely independent of each other.
The latter allows to build on top of currently unshipped features.

`

const parentBranchPromptTemplate = "Please specify the parent branch of %q:"

const perennialBranchOption = "<none> (perennial branch)"

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

func printParentBranchHeader(repo *git.ProdRepo) {
	cli.Printf(parentBranchHeaderTemplate, repo.Config.MainBranch())
}
