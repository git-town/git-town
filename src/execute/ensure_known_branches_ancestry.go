package execute

import (
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/validate"
)

// EnsureKnownBranchesAncestry makes sure the entire repo lineage is known.
// If needed, it queries the user for missing information.
// It returns the updated version of all information that is derived from the lineage.
//
// The purpose of this function is to implement proper cache invalidation.
// It ensures that all information derived from lineage gets updated when the lineage is updated.
func EnsureKnownBranchesAncestry(args EnsureKnownBranchesAncestryArgs) error {
	updated, err := validate.KnowsBranchesAncestors(validate.KnowsBranchesAncestorsArgs{
		Backend:          &args.Runner.Backend,
		BranchesToVerify: args.BranchesToVerify,
		Config:           args.Config,
		DefaultChoice:    args.DefaultChoice,
		DialogTestInputs: args.DialogTestInputs,
		LocalBranches:    args.LocalBranches,
		MainBranch:       args.MainBranch,
	})
	if err != nil {
		return err
	}
	if updated {
		args.Runner.Config.Reload()
	}
	return nil
}

type EnsureKnownBranchesAncestryArgs struct {
	BranchesToVerify gitdomain.LocalBranchNames
	Config           *config.Config
	DefaultChoice    gitdomain.LocalBranchName
	DialogTestInputs *components.TestInputs
	LocalBranches    gitdomain.BranchInfos
	MainBranch       gitdomain.LocalBranchName
	Runner           *git.ProdRunner
}
