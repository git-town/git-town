package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
)

type ParentBranches struct {
	parentBranchHeaderShown bool
}

// EnsureKnowsParentBranches asserts that the entire ancestry for all given branches
// is known to Git Town.
// Missing ancestry information is queried from the user.
func (pbd *ParentBranches) EnsureKnowsParentBranches(branches []string, repo *git.ProdRepo) error {
	for _, branch := range branches {
		if repo.Config.IsMainBranch(branch) || repo.Config.IsPerennialBranch(branch) || repo.Config.HasParentBranch(branch) {
			continue
		}
		err := pbd.AskForBranchAncestry(branch, repo.Config.MainBranch(), repo)
		if err != nil {
			return err
		}
		if pbd.parentBranchHeaderShown {
			fmt.Println()
		}
	}
	return nil
}

// AskForBranchAncestry prompts the user for all unknown ancestors of the given branch.
func (pbd *ParentBranches) AskForBranchAncestry(branch, defaultBranch string, repo *git.ProdRepo) error {
	current := branch
	var err error
	for {
		parent := repo.Config.ParentBranch(current)
		if parent == "" { //nolint:nestif
			pbd.printParentBranchHeader(repo)
			parent, err = pbd.AskForBranchParent(current, defaultBranch, repo)
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
			err = repo.Config.SetParent(current, parent)
			if err != nil {
				return err
			}
		}
		if parent == repo.Config.MainBranch() || repo.Config.IsPerennialBranch(parent) {
			break
		}
		current = parent
	}
	return nil
}

// AskForBranchParent prompts the user for the parent of the given branch.
func (pbd *ParentBranches) AskForBranchParent(branch, defaultBranch string, repo *git.ProdRepo) (string, error) {
	choices, err := repo.Silent.LocalBranchesMainFirst()
	if err != nil {
		return "", err
	}
	filteredChoices := filterOutSelfAndDescendants(branch, choices, repo)
	return askForBranch(askForBranchOptions{
		branches:      append([]string{perennialBranchOption}, filteredChoices...),
		prompt:        fmt.Sprintf(parentBranchPromptTemplate, branch),
		defaultBranch: defaultBranch,
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

func (pbd *ParentBranches) printParentBranchHeader(repo *git.ProdRepo) {
	if !pbd.parentBranchHeaderShown {
		pbd.parentBranchHeaderShown = true
		cli.Printf(parentBranchHeaderTemplate, repo.Config.MainBranch())
	}
}
