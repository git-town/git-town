package validate

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
)

// KnowsBranchesAncestry asserts that the entire ancestry for all given branches
// is known to Git Town.
// Missing ancestry information is queried from the user.
func KnowsBranchesAncestry(branches []string, repo *git.InternalCommands) error {
	mainBranch := repo.Config.MainBranch()
	for _, branch := range branches {
		err := KnowsBranchAncestry(branch, mainBranch, repo)
		if err != nil {
			return err
		}
	}
	return nil
}

// KnowsBranchAncestry prompts the user for all unknown ancestors of the given branch.
func KnowsBranchAncestry(branch, defaultBranch string, repo *git.InternalCommands) (err error) { //nolint:nonamedreturns // return value names are useful here
	headerShown := false
	currentBranch := branch
	if repo.Config.IsMainBranch(branch) || repo.Config.IsPerennialBranch(branch) || repo.Config.HasParentBranch(branch) {
		return nil
	}
	for {
		parent := repo.Config.ParentBranch(currentBranch)
		if parent == "" { //nolint:nestif
			if !headerShown {
				printParentBranchHeader(repo)
				headerShown = true
			}
			parent, err = EnterParent(currentBranch, defaultBranch, repo)
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

func printParentBranchHeader(repo *git.InternalCommands) {
	cli.Printf(parentBranchHeaderTemplate, repo.Config.MainBranch())
}

const parentBranchHeaderTemplate string = `
Feature branches can be branched directly off
%s or from other feature branches.

The former allows to develop and ship features completely independent of each other.
The latter allows to build on top of currently unshipped features.

`
