package validate

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// KnowsBranchAncestors prompts the user for all unknown ancestors of the given branch.
func KnowsBranchAncestors(branch gitdomain.LocalBranchName, args KnowsBranchAncestorsArgs) (bool, error) {
	currentBranch := branch
	if !args.Config.IsFeatureBranch(branch) {
		return false, nil
	}
	updated := false
	for {
		lineage := args.Backend.Config.Lineage
		parent, hasParent := lineage[currentBranch]
		if !hasParent { //nolint:nestif
			var aborted bool
			var err error
			parent, aborted, err = dialog.EnterParent(dialog.EnterParentArgs{
				Branch:          currentBranch,
				DialogTestInput: args.DialogTestInputs.Next(),
				LocalBranches:   args.LocalBranches,
				Lineage:         args.Config.Lineage,
				MainBranch:      args.MainBranch,
			})
			if err != nil {
				return false, err
			}
			if aborted {
				os.Exit(0)
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
	LocalBranches    gitdomain.LocalBranchNames
	Backend          *git.BackendCommands
	Config           *configdomain.FullConfig
	DialogTestInputs *dialog.TestInputs
	MainBranch       gitdomain.LocalBranchName
}

// KnowsBranchesAncestors asserts that the entire lineage for all given branches
// is known to Git Town.
// Prompts missing lineage information from the user.
// Indicates if the user made any changes to the ancestry.
func KnowsBranchesAncestors(args KnowsBranchesAncestorsArgs) (bool, error) {
	updated := false
	for _, branch := range args.LocalBranches {
		branchUpdated, err := KnowsBranchAncestors(branch.LocalName, KnowsBranchAncestorsArgs{
			MainBranch:       args.Config.MainBranch,
			Backend:          args.Backend,
			LocalBranches:    args.LocalBranches.Names(),
			Config:           args.Config,
			DialogTestInputs: args.DialogTestInputs,
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
	LocalBranches    gitdomain.BranchInfos
	Backend          *git.BackendCommands
	Config           *configdomain.FullConfig
	DialogTestInputs *dialog.TestInputs
}
