package execute

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/validate"
)

// LoadBranches loads the typically used information about Git branches using a single Git command.
func LoadBranches(args LoadBranchesArgs) (result LoadBranchesResult) {
	if args.HandleUnfinishedState {
		result.Exit, result.Err = validate.HandleUnfinishedState(&args.Repo.Runner, nil, args.Repo.RootDir)
		if result.Err != nil || result.Exit {
			return result
		}
	}
	if args.ValidateNoOpenChanges {
		var hasOpenChanges bool
		hasOpenChanges, result.Err = args.Repo.Runner.Backend.HasOpenChanges()
		if result.Err != nil {
			return result
		}
		result.Err = validate.NoOpenChanges(hasOpenChanges)
		if result.Err != nil {
			return result
		}
	}
	if args.Fetch {
		var remotes domain.Remotes
		remotes, result.Err = args.Repo.Runner.Backend.Remotes()
		if result.Err != nil {
			return result
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline {
			result.Err = args.Repo.Runner.Frontend.Fetch()
			if result.Err != nil {
				return result
			}
		}
	}
	var allBranches domain.BranchInfos
	allBranches, result.InitialBranch, result.Err = args.Repo.Runner.Backend.BranchInfos()
	if result.Err != nil {
		return result
	}
	result.BranchTypes = args.Repo.Runner.Config.BranchTypes()
	if args.ValidateIsConfigured {
		result.BranchTypes, result.Err = validate.IsConfigured(&args.Repo.Runner.Backend, allBranches, result.BranchTypes)
	}
	result.Snapshot = runstate.NewSnapshot(args.Repo.PartialSnapshot, allBranches)
	return result
}

type LoadBranchesArgs struct {
	Repo                  *OpenRepoResult
	Fetch                 bool
	HandleUnfinishedState bool
	ValidateIsConfigured  bool
	ValidateNoOpenChanges bool
}

type LoadBranchesResult struct {
	Snapshot      runstate.Snapshot
	InitialBranch domain.LocalBranchName
	BranchTypes   domain.BranchTypes
	Exit          bool
	Err           error
}
