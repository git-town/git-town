package execute

import (
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/git"
	"github.com/git-town/git-town/v10/src/validate"
)

// EnsureKnownBranchAncestry makes sure the lineage for the given branch is known.
// If needed, it queries the user for missing information.
// It returns the updated version of all information that is derived from the lineage.
//
// The purpose of this function is to implement proper cache invalidation.
// It ensures that all information derived from lineage gets updated when the lineage is updated.
func EnsureKnownBranchAncestry(branch domain.LocalBranchName, args EnsureKnownBranchAncestryArgs) (domain.BranchTypes, config.Lineage, error) {
	updated, err := validate.KnowsBranchAncestors(branch, validate.KnowsBranchAncestorsArgs{
		AllBranches:   args.AllBranches,
		Backend:       &args.Runner.Backend,
		BranchTypes:   args.BranchTypes,
		DefaultBranch: args.DefaultBranch,
		MainBranch:    args.MainBranch,
	})
	if err != nil {
		return args.BranchTypes, args.Lineage, err
	}
	if updated {
		args.Runner.Config.Reload()
		args.Lineage = args.Runner.Config.Lineage(args.Runner.Backend.Config.RemoveLocalConfigValue) // reload after ancestry change
		args.BranchTypes = args.Runner.Config.BranchTypes()
	}
	return args.BranchTypes, args.Lineage, nil
}

type EnsureKnownBranchAncestryArgs struct {
	AllBranches   domain.BranchInfos
	BranchTypes   domain.BranchTypes
	DefaultBranch domain.LocalBranchName
	Lineage       config.Lineage
	MainBranch    domain.LocalBranchName
	Runner        *git.ProdRunner
}
