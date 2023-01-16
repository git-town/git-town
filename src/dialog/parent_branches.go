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
func (pbd *ParentBranches) EnsureKnowsParentBranches(branchNames []string, repo *git.ProdRepo) error {
	for _, branchName := range branchNames {
		if repo.Config.IsMainBranch(branchName) || repo.Config.IsPerennialBranch(branchName) || repo.Config.HasParentBranch(branchName) {
			continue
		}
		err := pbd.AskForBranchAncestry(branchName, repo.Config.MainBranch(), repo)
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
func (pbd *ParentBranches) AskForBranchAncestry(branchName, defaultBranchName string, repo *git.ProdRepo) error {
	current := branchName
	var err error
	for {
		parent := repo.Config.ParentBranch(current)
		if parent == "" { //nolint:nestif
			pbd.printParentBranchHeader(repo)
			parent, err = pbd.AskForBranchParent(current, defaultBranchName, repo)
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
		if parent == repo.Config.MainBranch() || repo.Config.IsPerennialBranch(parent) {
			break
		}
		current = parent
	}
	return nil
}

// AskForBranchParent prompts the user for the parent of the given branch.
func (pbd *ParentBranches) AskForBranchParent(branchName, defaultBranchName string, repo *git.ProdRepo) (string, error) {
	choices, err := repo.Silent.LocalBranchesMainFirst()
	if err != nil {
		return "", err
	}
	filteredChoices := filterOutSelfAndDescendants(branchName, choices, repo)
	return askForBranch(askForBranchOptions{
		branchNames:       append([]string{perennialBranchOption}, filteredChoices...),
		prompt:            fmt.Sprintf(parentBranchPromptTemplate, branchName),
		defaultBranchName: defaultBranchName,
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

func filterOutSelfAndDescendants(branchName string, choices []string, repo *git.ProdRepo) []string {
	result := []string{}
	for _, choice := range choices {
		if choice == branchName || repo.Config.IsAncestorBranch(choice, branchName) {
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
