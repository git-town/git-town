package validate

import (
	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
)

// KnowsBranchesAncestors asserts that the entire lineage for all given branches
// is known to Git Town.
// Prompts missing lineage information from the user.
func KnowsBranchesAncestors(branches git.BranchesSyncStatus, mainBranch string, backend *git.BackendCommands, lineage config.Lineage, branchDurations config.BranchDurations) error {
	for _, branch := range branches {
		err := KnowsBranchAncestors(branch.Name, KnowsBranchAncestorsArgs{
			DefaultBranch:   mainBranch,
			Backend:         backend,
			AllBranches:     branches,
			Lineage:         lineage,
			BranchDurations: branchDurations,
			MainBranch:      mainBranch,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// KnowsBranchAncestors prompts the user for all unknown ancestors of the given branch.
// TODO: inject all dependencies.
func KnowsBranchAncestors(branch string, args KnowsBranchAncestorsArgs) (err error) {
	headerShown := false
	currentBranch := branch
	if args.BranchDurations.IsMainBranch(branch) || args.BranchDurations.IsPerennialBranch(branch) {
		return nil
	}
	for {
		parent := args.Backend.Config.Lineage()[currentBranch] // need to reload the lineage here because ancestry data was changed
		if parent == "" {                                      //nolint:nestif
			if !headerShown {
				printParentBranchHeader(args.MainBranch)
				headerShown = true
			}
			parent, err = EnterParent(currentBranch, args.DefaultBranch, args.Lineage, args.AllBranches)
			if err != nil {
				return
			}
			if parent == perennialBranchOption {
				err = args.Backend.Config.AddToPerennialBranches(currentBranch)
				if err != nil {
					return
				}
				break
			}
			err = args.Backend.Config.SetParent(currentBranch, parent)
			if err != nil {
				return
			}
		}
		if !args.BranchDurations.IsFeatureBranch(parent) {
			break
		}
		currentBranch = parent
	}
	return
}

type KnowsBranchAncestorsArgs struct {
	AllBranches     git.BranchesSyncStatus
	Backend         *git.BackendCommands
	BranchDurations config.BranchDurations
	DefaultBranch   string
	Lineage         config.Lineage
	MainBranch      string
}

func printParentBranchHeader(mainBranch string) {
	cli.Printf(parentBranchHeaderTemplate, mainBranch)
}

const parentBranchHeaderTemplate string = `
Feature branches can be branched directly off
%s or from other feature branches.

The former allows to develop and ship features completely independent of each other.
The latter allows to build on top of currently unshipped features.

`
