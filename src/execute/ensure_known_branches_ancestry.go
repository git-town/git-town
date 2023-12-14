package execute

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/validate"
)

// EnsureKnownBranchesAncestry makes sure the entire repo lineage is known.
// If needed, it queries the user for missing information.
// It returns the updated version of all information that is derived from the lineage.
//
// The purpose of this function is to implement proper cache invalidation.
// It ensures that all information derived from lineage gets updated when the lineage is updated.
func EnsureKnownBranchesAncestry(args EnsureKnownBranchesAncestryArgs) (domain.BranchTypes, configdomain.Lineage, error) {
	updated, err := validate.KnowsBranchesAncestors(validate.KnowsBranchesAncestorsArgs{
		AllBranches: args.AllBranches,
		Backend:     &args.Runner.Backend,
		BranchTypes: args.BranchTypes,
		MainBranch:  args.MainBranch,
	})
	if err != nil {
		return args.BranchTypes, args.Lineage, err
	}
	if updated {
		args.Runner.GitTown.Reload()
		args.Lineage = args.Runner.GitTown.Lineage(args.Runner.Backend.GitTown.RemoveLocalConfigValue) // reload after ancestry change
		args.BranchTypes = args.Runner.GitTown.BranchTypes()
	}
	return args.BranchTypes, args.Lineage, nil
}

type EnsureKnownBranchesAncestryArgs struct {
	AllBranches domain.BranchInfos
	BranchTypes domain.BranchTypes
	Lineage     configdomain.Lineage
	MainBranch  domain.LocalBranchName
	Runner      *git.ProdRunner
}
