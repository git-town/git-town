package validate

import (
	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/cli/io"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git"
)

// KnowsBranchAncestors prompts the user for all unknown ancestors of the given branch.
func KnowsBranchAncestors(branch domain.LocalBranchName, args KnowsBranchAncestorsArgs) (bool, error) {
	headerShown := false
	currentBranch := branch
	if !args.BranchTypes.IsFeatureBranch(branch) {
		return false, nil
	}
	updated := false
	for {
		lineage := args.Backend.GitTown.Lineage(args.Backend.GitTown.RemoveLocalConfigValue) // need to load a fresh lineage here because ancestry data was changed
		parent, hasParent := lineage[currentBranch]
		var err error
		if !hasParent { //nolint:nestif
			if !headerShown {
				printParentBranchHeader(args.MainBranch)
				headerShown = true
			}
			parent, err = dialog.EnterParent(currentBranch, args.DefaultBranch, lineage, args.AllBranches)
			if err != nil {
				return false, err
			}
			if parent.String() == dialog.PerennialBranchOption {
				err = args.Backend.GitTown.AddToPerennialBranches(currentBranch)
				if err != nil {
					return false, err
				}
				updated = true
				break
			}
			err = args.Backend.GitTown.SetParent(currentBranch, parent)
			if err != nil {
				return false, err
			}
			updated = true
		}
		if !args.BranchTypes.IsFeatureBranch(parent) {
			break
		}
		currentBranch = parent
	}
	return updated, nil
}

type KnowsBranchAncestorsArgs struct {
	AllBranches   domain.BranchInfos
	Backend       *git.BackendCommands
	BranchTypes   domain.BranchTypes
	DefaultBranch domain.LocalBranchName
	MainBranch    domain.LocalBranchName
}

// KnowsBranchesAncestors asserts that the entire lineage for all given branches
// is known to Git Town.
// Prompts missing lineage information from the user.
// Indicates if the user made any changes to the ancestry.
func KnowsBranchesAncestors(args KnowsBranchesAncestorsArgs) (bool, error) {
	updated := false
	for _, branch := range args.AllBranches {
		branchUpdated, err := KnowsBranchAncestors(branch.LocalName, KnowsBranchAncestorsArgs{
			DefaultBranch: args.MainBranch,
			Backend:       args.Backend,
			AllBranches:   args.AllBranches,
			BranchTypes:   args.BranchTypes,
			MainBranch:    args.MainBranch,
		})
		if err != nil {
			return updated, err
		}
		if branchUpdated {
			updated = true
		}
	}
	return updated, nil
}

type KnowsBranchesAncestorsArgs struct {
	AllBranches domain.BranchInfos
	Backend     *git.BackendCommands
	BranchTypes domain.BranchTypes
	MainBranch  domain.LocalBranchName
}

func printParentBranchHeader(mainBranch domain.LocalBranchName) {
	io.Printf(parentBranchHeaderTemplate, mainBranch)
}

const parentBranchHeaderTemplate string = `
Feature branches can be branched directly off
%s or from other feature branches.

The former allows to develop and ship features completely independent of each other.
The latter allows to build on top of currently unshipped features.

`
