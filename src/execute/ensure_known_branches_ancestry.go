package execute

import (
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
func EnsureKnownBranchesAncestry(args EnsureKnownBranchesAncestryArgs) (configdomain.BranchTypes, configdomain.Lineage, error) {
	updated, err := validate.KnowsBranchesAncestors(validate.KnowsBranchesAncestorsArgs{
		AllBranches: args.AllBranches,
		Backend:     &args.Runner.Backend,
		Config:      args.Config,
	})
	if err != nil {
		return args.BranchTypes, args.Config.Lineage, err
	}
	if updated {
		args.Runner.Config.Reload()
		args.Config.Lineage = args.Runner.Config.Lineage // reload after ancestry change
		args.BranchTypes = args.Runner.Config.BranchTypes()
	}
	return args.BranchTypes, args.Config.Lineage, nil
}

type EnsureKnownBranchesAncestryArgs struct {
	Config      *configdomain.FullConfig
	AllBranches gitdomain.BranchInfos
	BranchTypes configdomain.BranchTypes
	Runner      *git.ProdRunner
}
