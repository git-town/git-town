package execute

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/validate"
)

// LoadBranches loads the typically used information about Git branches using a single Git command.
func LoadBranches(args LoadBranchesArgs) (snapshot runstate.Snapshot, exit bool, err error) {
	if args.HandleUnfinishedState {
		exit, err := validate.HandleUnfinishedState(&args.Repo.Runner, nil, args.Repo.RootDir)
		if err != nil || exit {
			return snapshot, exit, err
		}
	}
	if args.ValidateNoOpenChanges {
		hasOpenChanges, err := args.Repo.Runner.Backend.HasOpenChanges()
		if err != nil {
			return snapshot, false, err
		}
		err = validate.NoOpenChanges(hasOpenChanges)
		if err != nil {
			return snapshot, false, err
		}
	}
	if args.Fetch {
		var remotes domain.Remotes
		remotes, err := args.Repo.Runner.Backend.Remotes()
		if err != nil {
			return snapshot, false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline {
			err = args.Repo.Runner.Frontend.Fetch()
			if err != nil {
				return snapshot, false, err
			}
		}
	}
	allBranches, initialBranch, err := args.Repo.Runner.Backend.BranchInfos()
	if err != nil {
		return snapshot, false, err
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
	Repo                  *OpenRepoResult
	Fetch                 bool
	HandleUnfinishedState bool
	ValidateIsConfigured  bool
	ValidateNoOpenChanges bool
}
