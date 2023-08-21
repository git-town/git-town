package execute

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/validate"
)

// LoadBranches loads the typically used information about Git branches using a single Git command.
func LoadBranches(args LoadBranchesArgs) (domain.Branches, bool, error) {
	if args.HandleUnfinishedState {
		exit, err := validate.HandleUnfinishedState(&args.Repo.Runner, nil, args.Repo.RootDir)
		if err != nil || exit {
			return domain.EmptyBranches(), exit, err
		}
	}
	allBranches, initialBranch, err := args.Repo.Runner.Backend.BranchInfos()
	if err != nil {
		return domain.EmptyBranches(), false, err
	}
	branchTypes := args.Repo.Runner.Config.BranchTypes()
	result := domain.Branches{
		All:     allBranches,
		Types:   branchTypes,
		Initial: initialBranch,
	}
	if args.ValidateIsConfigured {
		result.Types, err = validate.IsConfigured(&args.Repo.Runner.Backend, result)
	}
	return result, false, err
}

type LoadBranchesArgs struct {
	Repo                  *RepoData
	HandleUnfinishedState bool
	ValidateIsConfigured  bool
}
