package validate

import (
	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/cli/io"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// KnowsBranchAncestors prompts the user for all unknown ancestors of the given branch.
func KnowsBranchAncestors(branch gitdomain.LocalBranchName, args KnowsBranchAncestorsArgs) (bool, error) {
	headerShown := false
	currentBranch := branch
	if !args.Config.IsFeatureBranch(branch) {
		return false, nil
	}
	updated := false
	for {
		lineage := args.Backend.Config.Lineage
		parent, hasParent := lineage[currentBranch]
		var err error
		if !hasParent { //nolint:nestif
			if !headerShown {
				printParentBranchHeader(args.Config.MainBranch)
				headerShown = true
			}
			parent, err = dialog.EnterParent(currentBranch, args.DefaultBranch, lineage, args.AllBranches)
			if err != nil {
				return false, err
			}
			if parent.String() == dialog.PerennialBranchOption {
				err = args.Backend.Config.AddToPerennialBranches(currentBranch)
				if err != nil {
					return false, err
				}
				updated = true
				break
			}
			err = args.Backend.Config.SetParent(currentBranch, parent)
			if err != nil {
				return false, err
			}
			updated = true
		}
		if !args.Config.IsFeatureBranch(parent) {
			break
		}
		currentBranch = parent
	}
	return updated, nil
}

type KnowsBranchAncestorsArgs struct {
	AllBranches   gitdomain.BranchInfos
	Backend       *git.BackendCommands
	Config        *configdomain.FullConfig
	DefaultBranch gitdomain.LocalBranchName
}

// KnowsBranchesAncestors asserts that the entire lineage for all given branches
// is known to Git Town.
// Prompts missing lineage information from the user.
// Indicates if the user made any changes to the ancestry.
func KnowsBranchesAncestors(args KnowsBranchesAncestorsArgs) (bool, error) {
	updated := false
	for _, branch := range args.AllBranches {
		branchUpdated, err := KnowsBranchAncestors(branch.LocalName, KnowsBranchAncestorsArgs{
			DefaultBranch: args.Config.MainBranch,
			Backend:       args.Backend,
			AllBranches:   args.AllBranches,
			Config:        args.Config,
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
	AllBranches gitdomain.BranchInfos
	Backend     *git.BackendCommands
	Config      *configdomain.FullConfig
}

func printParentBranchHeader(mainBranch gitdomain.LocalBranchName) {
	io.Printf(parentBranchHeaderTemplate, mainBranch)
}

const parentBranchHeaderTemplate string = `
Feature branches can be branched directly off
%s or from other feature branches.

The former allows to develop and ship features completely independent of each other.
The latter allows to build on top of currently unshipped features.

`
