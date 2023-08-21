package execute

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/validate"
)

// LoadBranches loads the typically used information about Git branches using a single Git command.
func LoadBranches(args LoadBranchesArgs) (domain.Branches, error) {
	allBranches, initialBranch, err := args.Runner.Backend.BranchInfos()
	if err != nil {
		return domain.EmptyBranches(), err
	}
	branchTypes := args.Runner.Config.BranchTypes()
	result := domain.Branches{
		All:     allBranches,
		Types:   branchTypes,
		Initial: initialBranch,
	}
	if args.ValidateIsConfigured {
		result.Types, err = validate.IsConfigured(&args.Runner.Backend, result)
	}
	return result, err
}

type LoadBranchesArgs struct {
	Runner               *git.ProdRunner
	ValidateIsConfigured bool
}
