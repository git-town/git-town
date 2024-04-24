package validate

import (
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

// KnowsBranchAncestors prompts the user for all unknown ancestors of the given branch.
func KnowsBranchAncestors(branch gitdomain.LocalBranchName, args KnowsBranchAncestorsArgs) (bool, error) {
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

type KnowsBranchAncestorsArgs struct {
	Backend          *git.BackendCommands
	Config           *config.Config
	DefaultChoice    gitdomain.LocalBranchName
	DialogTestInputs *components.TestInputs
	LocalBranches    gitdomain.LocalBranchNames
	MainBranch       gitdomain.LocalBranchName
}

// KnowsBranchesAncestors asserts that the entire lineage for all given branches
// is known to Git Town.
// Prompts missing lineage information from the user.
// Indicates if the user made any changes to the ancestry.
func KnowsBranchesAncestors(args KnowsBranchesAncestorsArgs) (bool, error) {
	updated := false
	for _, branch := range args.BranchesToVerify {
		branchUpdated, err := KnowsBranchAncestors(branch, KnowsBranchAncestorsArgs{
			Backend:          args.Backend,
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
	Config           *config.Config
	DefaultChoice    gitdomain.LocalBranchName
	DialogTestInputs *components.TestInputs
	LocalBranches    gitdomain.BranchInfos
	MainBranch       gitdomain.LocalBranchName
}
