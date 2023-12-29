package execute

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/validate"
)

// LoadBranches loads the typically used information about Git branches using a single Git command.
func LoadBranches(args LoadBranchesArgs) (configdomain.Branches, gitdomain.BranchesStatus, gitdomain.StashSize, bool, error) {
	var branchesSnapshot gitdomain.BranchesStatus
	var err error
	stashSnapshot, err := args.Repo.Runner.Backend.StashSize()
	if err != nil {
		return configdomain.EmptyBranches(), gitdomain.EmptyBranchesSnapshot(), stashSnapshot, false, err
	}
	if args.HandleUnfinishedState {
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return configdomain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
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
			Run:                     args.Repo.Runner,
		})
		if err != nil || exit {
			return configdomain.EmptyBranches(), branchesSnapshot, stashSnapshot, exit, err
		}
	}
	if args.ValidateNoOpenChanges {
		repoStatus, err := args.Repo.Runner.Backend.RepoStatus()
		if err != nil {
			return configdomain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
		}
		err = validate.NoOpenChanges(repoStatus.OpenChanges)
		if err != nil {
			return configdomain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
		}
	}
	if args.Fetch {
		var remotes gitdomain.Remotes
		remotes, err := args.Repo.Runner.Backend.Remotes()
		if err != nil {
			return configdomain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline.Bool() {
			err = args.Repo.Runner.Frontend.Fetch()
			if err != nil {
				return configdomain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
			}
		}
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return configdomain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
		}
	}
	if branchesSnapshot.IsEmpty() {
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return configdomain.EmptyBranches(), branchesSnapshot, stashSnapshot, false, err
		}
	}
	branchTypes := args.Repo.Runner.Config.BranchTypes()
	branches := configdomain.Branches{
		All:     branchesSnapshot.Branches,
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
	Lineage               configdomain.Lineage
	PushHook              configdomain.PushHook
	ValidateIsConfigured  bool
	ValidateNoOpenChanges bool
}
