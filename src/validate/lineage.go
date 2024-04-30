package validate

import (
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

// Lineage validates that the given lineage contains the ancestry for all given branches.
// Prompts missing lineage information from the user and updates persisted lineage as needed.
// Returns the validated Lineage.
func Lineage(args KnowsBranchesAncestorsArgs) (configdomain.Lineage, error) {
	// step 1: determine all branches for which the parent must be known
	// step 2: for each branch: check the ancestor
	// step 3: if missing: ask user and add the ancestry info to the validated lineage
	// step 4: add the parent to the list of branches that need to be verified

	updated := false
	for _, branch := range args.BranchesToVerify {
		branchUpdated, err := knowsBranchAncestors(branch, knowsBranchAncestorsArgs{
			Config:           args.Config,
			DefaultChoice:    args.DefaultChoice,
			DialogTestInputs: args.DialogTestInputs,
			LocalBranches:    args.LocalBranches.Names(),
			MainBranch:       args.MainBranch,
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
	Backend          *git.BackendCommands
	BranchesToVerify gitdomain.LocalBranchNames
	Config           *config.UnvalidatedConfig
	DefaultChoice    gitdomain.LocalBranchName
	DialogTestInputs *components.TestInputs
	LocalBranches    gitdomain.BranchInfos
	MainBranch       gitdomain.LocalBranchName
}

// knowsBranchAncestors prompts the user for all unknown ancestors of the given branch.
func knowsBranchAncestors(branch gitdomain.LocalBranchName, args knowsBranchAncestorsArgs) (bool, error) {
	currentBranch := branch
	if args.Config.FullConfig.IsMainOrPerennialBranch(branch) || args.Config.FullConfig.IsObservedBranch(branch) || args.Config.FullConfig.IsContributionBranch(branch) {
		return false, nil
	}
	updated := false
	for {
		lineage := args.Config.FullConfig.Lineage
		parent, hasParent := lineage[currentBranch]
		if !hasParent { //nolint:nestif
			var err error
			outcome, selectedBranch, err := dialog.Parent(dialog.ParentArgs{
				Branch:          currentBranch,
				DefaultChoice:   args.DefaultChoice,
				DialogTestInput: args.DialogTestInputs.Next(),
				Lineage:         args.Config.FullConfig.Lineage,
				LocalBranches:   args.LocalBranches,
				MainBranch:      args.MainBranch,
			})
			if err != nil {
				return false, err
			}
			switch outcome {
			case dialog.ParentOutcomeAborted:
				os.Exit(0)
			case dialog.ParentOutcomePerennialBranch:
				err = args.Config.AddToPerennialBranches(currentBranch)
				if err != nil {
					return false, err
				}
				updated = true
				return updated, nil
			case dialog.ParentOutcomeSelectedParent:
				parent = selectedBranch
				err = args.Config.SetParent(currentBranch, parent)
				if err != nil {
					return false, err
				}
				updated = true
			}
		}
		if args.Config.FullConfig.IsMainOrPerennialBranch(parent) {
			break
		}
		currentBranch = parent
	}
	return updated, nil
}

type knowsBranchAncestorsArgs struct {
	Config           *config.Config
	DefaultChoice    gitdomain.LocalBranchName
	DialogTestInputs *components.TestInputs
	LocalBranches    gitdomain.LocalBranchNames
	MainBranch       gitdomain.LocalBranchName
}
