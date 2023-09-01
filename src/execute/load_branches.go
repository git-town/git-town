package execute

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/validate"
)

// LoadSnapshot loads the typically used information about Git branches using a single Git command.
func LoadSnapshot(args LoadBranchesArgs) (runstate.Snapshot, domain.LocalBranchName, bool, error) {
	if args.HandleUnfinishedState {
		exit, err := validate.HandleUnfinishedState(&args.Repo.Runner, nil, args.Repo.RootDir)
		if err != nil || exit {
			return runstate.EmptySnapshot(), domain.LocalBranchName{}, exit, err
		}
	}
	if args.ValidateNoOpenChanges {
		hasOpenChanges, err := args.Repo.Runner.Backend.HasOpenChanges()
		if err != nil {
			return runstate.EmptySnapshot(), domain.LocalBranchName{}, false, err
		}
		err = validate.NoOpenChanges(hasOpenChanges)
		if err != nil {
			return runstate.EmptySnapshot(), domain.LocalBranchName{}, false, err
		}
	}
	if args.Fetch {
		var remotes domain.Remotes
		remotes, err := args.Repo.Runner.Backend.Remotes()
		if err != nil {
			return runstate.EmptySnapshot(), domain.LocalBranchName{}, false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline {
			err = args.Repo.Runner.Frontend.Fetch()
			if err != nil {
				return runstate.EmptySnapshot(), domain.LocalBranchName{}, false, err
			}
		}
	}
	allBranches, initialBranch, err := args.Repo.Runner.Backend.BranchInfos()
	if err != nil {
		return runstate.EmptySnapshot(), domain.LocalBranchName{}, false, err
	}
	branchTypes := args.Repo.Runner.Config.BranchTypes()
	if args.ValidateIsConfigured {
		branchTypes, err = validate.IsConfigured(&args.Repo.Runner.Backend, allBranches, branchTypes)
	}
	return runstate.Snapshot{
		PartialSnapshot: runstate.NewPartialSnapshot(args.Repo.Runner.Config.Git, args.Repo.),
		Branches:        allBranches,
	}, initialBranch, false, err
}

type LoadBranchesArgs struct {
	Repo                  *OpenRepoResult
	Fetch                 bool
	HandleUnfinishedState bool
	ValidateIsConfigured  bool
	ValidateNoOpenChanges bool
}
