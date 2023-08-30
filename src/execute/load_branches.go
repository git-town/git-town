package execute

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/validate"
)

// LoadSnapshot loads the typically used information about Git branches using a single Git command.
func LoadSnapshot(args LoadBranchesArgs) (runstate.Snapshot, bool, error) {
	if args.HandleUnfinishedState {
		exit, err := validate.HandleUnfinishedState(&args.Repo.Runner, nil, args.Repo.RootDir)
		if err != nil || exit {
			return runstate.EmptySnapshot(), exit, err
		}
	}
	if args.ValidateNoOpenChanges {
		hasOpenChanges, err := args.Repo.Runner.Backend.HasOpenChanges()
		if err != nil {
			return runstate.EmptySnapshot(), false, err
		}
		err = validate.NoOpenChanges(hasOpenChanges)
		if err != nil {
			return runstate.EmptySnapshot(), false, err
		}
	}
	if args.Fetch {
		var remotes domain.Remotes
		remotes, err := args.Repo.Runner.Backend.Remotes()
		if err != nil {
			return runstate.EmptySnapshot(), false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline {
			err = args.Repo.Runner.Frontend.Fetch()
			if err != nil {
				return runstate.EmptySnapshot(), false, err
			}
		}
	}
	allBranches, initialBranch, err := args.Repo.Runner.Backend.BranchInfos()
	if err != nil {
		return runstate.EmptySnapshot(), false, err
	}
	branchTypes := args.Repo.Runner.Config.BranchTypes()
	branches := domain.Branches{
		All:     allBranches,
		Types:   branchTypes,
		Initial: initialBranch,
	}
	if args.ValidateIsConfigured {
		branches.Types, err = validate.IsConfigured(&args.Repo.Runner.Backend, branches)
	}
	return runstate.Snapshot{
		PartialSnapshot: runstate.NewPartialSnapshot(args.Repo.Runner.Config.Git),
		Branches:        branches.All,
	}, false, err
}

type LoadBranchesArgs struct {
	Repo                  *OpenRepoResult
	Fetch                 bool
	HandleUnfinishedState bool
	ValidateIsConfigured  bool
	ValidateNoOpenChanges bool
}
