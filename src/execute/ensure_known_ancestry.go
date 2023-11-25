package execute

import (
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/git"
	"github.com/git-town/git-town/v10/src/validate"
)

// EnsureKnownAncestry makes sure the entire branch lineage is known.
// If needed, it queries the user for missing information.
// It returns the updated version of all information that is derived from the lineage.
//
// The purpose of this function is to implement proper cache invalidation.
// It ensures that all Git Town commands use the updated information when the lineage is updated.
func EnsureKnowsBranchAncestry(branch domain.LocalBranchName, args EnsureKnowsBranchAncestryArgs) (domain.BranchTypes, config.Lineage, error) {
	knowsBranchAncestorArgs := validate.KnowsBranchAncestorsArgs{
		AllBranches:   args.AllBranches,
		Backend:       &args.Runner.Backend,
		BranchTypes:   args.BranchTypes,
		DefaultBranch: args.DefaultBranch,
		MainBranch:    args.MainBranch,
	}
	updated, err := validate.KnowsBranchAncestors(branch, knowsBranchAncestorArgs)
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

type EnsureKnowsBranchAncestryArgs struct {
	AllBranches   domain.BranchInfos
	Runner        *git.ProdRunner
	BranchTypes   domain.BranchTypes
	DefaultBranch domain.LocalBranchName
	Lineage       config.Lineage
	MainBranch    domain.LocalBranchName
}
