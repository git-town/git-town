package validate

import (
	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// KnowsBranchesAncestors asserts that the entire lineage for all given branches
// is known to Git Town.
// Prompts missing lineage information from the user.
// Indicates if the user made any changes to the ancestry.
func KnowsBranchesAncestors(args KnowsBranchesAncestorsArgs) (bool, error) {
	updated := false
	for _, branch := range args.AllBranches {
		branchUpdated, err := KnowsBranchAncestors(branch.Name, KnowsBranchAncestorsArgs{
			DefaultBranch:   args.MainBranch,
			Backend:         args.Backend,
			AllBranches:     args.AllBranches,
			Lineage:         args.Lineage,
			BranchDurations: args.BranchDurations,
			MainBranch:      args.MainBranch,
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
	AllBranches     domain.BranchInfos
	Backend         *git.BackendCommands
	BranchDurations domain.BranchDurations
	Lineage         config.Lineage
	MainBranch      domain.LocalBranchName
}

// KnowsBranchAncestors prompts the user for all unknown ancestors of the given branch.
func KnowsBranchAncestors(branch domain.LocalBranchName, args KnowsBranchAncestorsArgs) (bool, error) {
	headerShown := false
	currentBranch := branch
	if !args.BranchDurations.IsFeatureBranch(branch) {
		return false, nil
	}
	updated := false
	for {
		// TODO: reload the lineage at the end of the loop
		parent, hasParent := args.Backend.Config.Lineage()[currentBranch] // need to reload the lineage here because ancestry data was changed
		var err error
		if !hasParent { //nolint:nestif
			if !headerShown {
				printParentBranchHeader(args.MainBranch)
				headerShown = true
			}
			parent, err = dialog.EnterParent(currentBranch, args.DefaultBranch, args.Lineage, args.AllBranches)
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
		if !args.BranchDurations.IsFeatureBranch(parent) {
			break
		}
		currentBranch = parent
	}
	return updated, nil
}

type KnowsBranchAncestorsArgs struct {
	AllBranches     domain.BranchInfos
	Backend         *git.BackendCommands
	BranchDurations domain.BranchDurations
	DefaultBranch   domain.LocalBranchName
	Lineage         config.Lineage
	MainBranch      domain.LocalBranchName
}

func printParentBranchHeader(mainBranch domain.LocalBranchName) {
	cli.Printf(parentBranchHeaderTemplate, mainBranch)
}

const parentBranchHeaderTemplate string = `
Feature branches can be branched directly off
%s or from other feature branches.

The former allows to develop and ship features completely independent of each other.
The latter allows to build on top of currently unshipped features.

`
