package execute

import (
	"github.com/git-town/git-town/v9/src/config"
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
	if args.ValidateNoOpenChanges {
		hasOpenChanges, err := args.Repo.Runner.Backend.HasOpenChanges()
		if err != nil {
			return domain.EmptyBranches(), false, err
		}
		err = validate.NoOpenChanges(hasOpenChanges)
		if err != nil {
			return domain.EmptyBranches(), false, err
		}
	}
	if args.Fetch {
		var remotes config.Remotes
		remotes, err := args.Repo.Runner.Backend.Remotes()
		if err != nil {
			return domain.EmptyBranches(), false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline {
			err = args.Repo.Runner.Frontend.Fetch()
			if err != nil {
				return domain.EmptyBranches(), false, err
			}
		}
	}
	allBranches, initialBranch, err := args.Repo.Runner.Backend.BranchInfos()
	if err != nil {
		return domain.EmptyBranches(), false, err
	}
	if args.HandleUnfinishedState {
		exit, err := validate.HandleUnfinishedState(&args.Repo.Runner, nil, args.Repo.RootDir, allBranches)
		if err != nil || exit {
			return domain.EmptyBranches(), exit, err
		}
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
	Fetch                 bool
	HandleUnfinishedState bool
	ValidateIsConfigured  bool
	ValidateNoOpenChanges bool
}
