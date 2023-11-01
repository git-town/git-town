package execute

import (
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/validate"
)

// LoadBranches loads the typically used information about Git branches using a single Git command.
func LoadBranches(args LoadBranchesArgs) (domain.Branches, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
	var branchesSnapshot domain.BranchesSnapshot
	var err error
	stashSnapshot, err := args.Repo.Runner.Backend.StashSnapshot()
	if err != nil {
		return domain.EmptyBranches(), domain.EmptyBranchesSnapshot(), stashSnapshot, false, err
	}
	if args.HandleUnfinishedState {
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return domain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
		}
		exit, err := validate.HandleUnfinishedState(validate.UnfinishedStateArgs{
			Connector:               nil,
			Verboe:                  args.Verbose,
			InitialBranchesSnapshot: branchesSnapshot,
			InitialConfigSnapshot:   args.Repo.ConfigSnapshot,
			InitialStashSnapshot:    stashSnapshot,
			Lineage:                 args.Lineage,
			PushHook:                args.PushHook,
			RootDir:                 args.Repo.RootDir,
			Run:                     &args.Repo.Runner,
		})
		if err != nil || exit {
			return domain.EmptyBranches(), branchesSnapshot, stashSnapshot, exit, err
		}
	}
	if args.ValidateNoOpenChanges {
		repoStatus, err := args.Repo.Runner.Backend.RepoStatus()
		if err != nil {
			return domain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
		}
		err = validate.NoOpenChanges(repoStatus.OpenChanges)
		if err != nil {
			return domain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
		}
	}
	if args.Fetch {
		var remotes domain.Remotes
		remotes, err := args.Repo.Runner.Backend.Remotes()
		if err != nil {
			return domain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline {
			err = args.Repo.Runner.Frontend.Fetch()
			if err != nil {
				return domain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
			}
		}
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return domain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
		}
	}
	if branchesSnapshot.IsEmpty() {
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return domain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
		}
	}
	branchTypes := args.Repo.Runner.Config.BranchTypes()
	branches := domain.Branches{
		All:     branchesSnapshot.Branches.Clone(),
		Types:   branchTypes,
		Initial: branchesSnapshot.Active,
	}
	if args.ValidateIsConfigured {
		branches.Types, err = validate.IsConfigured(&args.Repo.Runner.Backend, branches)
	}
	return branches, branchesSnapshot, stashSnapshot, false, err
}

type LoadBranchesArgs struct {
	Repo                  *OpenRepoResult
	Fetch                 bool
	Verbose               bool
	HandleUnfinishedState bool
	Lineage               config.Lineage
	PushHook              bool
	ValidateIsConfigured  bool
	ValidateNoOpenChanges bool
}
