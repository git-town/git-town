package validate

import (
	"github.com/git-town/git-town/v8/src/cli"
	"github.com/git-town/git-town/v8/src/git"
)

// KnowsBranchesAncestry asserts that the entire ancestry for all given branches
// is known to Git Town.
// Missing ancestry information is queried from the user.
func KnowsBranchesAncestry(branches []string, backend *git.BackendCommands) error {
	mainBranch := backend.MainBranch()
	for _, branch := range branches {
		err := KnowsBranchAncestry(branch, mainBranch, backend)
		if err != nil {
			return err
		}
	}
	return nil
}

// KnowsBranchAncestry prompts the user for all unknown ancestors of the given branch.
func KnowsBranchAncestry(branch, defaultBranch string, backend *git.BackendCommands) (err error) { //nolint:nonamedreturns // return value names are useful here
	headerShown := false
	currentBranch := branch
	if backend.IsMainBranch(branch) || backend.IsPerennialBranch(branch) || backend.HasParentBranch(branch) {
		return nil
	}
	for {
		parent := backend.ParentBranch(currentBranch)
		if parent == "" { //nolint:nestif
			if !headerShown {
				printParentBranchHeader(backend)
				headerShown = true
			}
			parent, err = EnterParent(currentBranch, defaultBranch, backend)
			if err != nil {
				return
			}
			if parent == perennialBranchOption {
				err = backend.AddToPerennialBranches(currentBranch)
				if err != nil {
					return
				}
				break
			}
			err = backend.SetParent(currentBranch, parent)
			if err != nil {
				return
			}
		}
		if parent == backend.MainBranch() || backend.IsPerennialBranch(parent) {
			break
		}
		currentBranch = parent
	}
	return
}

func printParentBranchHeader(backend *git.BackendCommands) {
	cli.Printf(parentBranchHeaderTemplate, backend.MainBranch())
}

const parentBranchHeaderTemplate string = `
Feature branches can be branched directly off
%s or from other feature branches.

The former allows to develop and ship features completely independent of each other.
The latter allows to build on top of currently unshipped features.

`
