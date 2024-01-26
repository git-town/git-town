package execute

import (
	"github.com/git-town/git-town/v11/src/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/validate"
)

// EnsureKnownBranchesAncestry makes sure the entire repo lineage is known.
// If needed, it queries the user for missing information.
// It returns the updated version of all information that is derived from the lineage.
//
// The purpose of this function is to implement proper cache invalidation.
// It ensures that all information derived from lineage gets updated when the lineage is updated.
func EnsureKnownBranchesAncestry(args EnsureKnownBranchesAncestryArgs) error {
	updated, err := validate.KnowsBranchesAncestors(validate.KnowsBranchesAncestorsArgs{
		LocalBranches:    args.LocalBranches,
		Backend:          &args.Runner.Backend,
		Config:           args.Config,
		DialogTestInputs: args.DialogTestInputs,
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
	Config           *configdomain.FullConfig
	LocalBranches    gitdomain.BranchInfos
	DialogTestInputs *dialogcomponents.TestInputs
	Runner           *git.ProdRunner
}
