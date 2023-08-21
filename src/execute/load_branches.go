package execute

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/validate"
)

// LoadBranches loads the typically used information about Git branches using a single Git command.
func LoadBranches(pr *git.ProdRunner, args LoadBranchesArgs) (domain.Branches, bool, error) {
	allBranches, initialBranch, err := pr.Backend.BranchInfos()
	if err != nil {
		return domain.EmptyBranches(), false, err
	}
	branchTypes := pr.Config.BranchTypes()
	result := domain.Branches{
		All:     allBranches,
		Types:   branchTypes,
		Initial: initialBranch,
	}
	if args.HandleUnfinishedState {
		exit, err := validate.HandleUnfinishedState(&pr, nil, rootDir)
		if err != nil || exit {
			return
		}
	}
	if args.ValidateIsConfigured {
		result.Types, err = validate.IsConfigured(&pr.Backend, result)
	}
	return result, err
}

type LoadBranchesArgs struct {
	HandleUnfinishedState bool
	ValidateIsConfigured  bool
}
