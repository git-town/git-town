package execute

import (
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/validate"
)

// LoadRepoSnapshot loads the initial snapshot of the Git repo.
func LoadRepoSnapshot(args LoadRepoSnapshotArgs) (gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	var branchesSnapshot gitdomain.BranchesSnapshot
	var err error
	stashSize, err := args.Repo.Runner.Backend.StashSize()
	if err != nil {
		return branchesSnapshot, stashSize, false, err
	}
	if args.HandleUnfinishedState {
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
		exit, err := validate.HandleUnfinishedState(validate.UnfinishedStateArgs{
			Connector:               nil,
			CurrentBranch:           branchesSnapshot.Active,
			DialogTestInputs:        args.DialogTestInputs,
			HasOpenChanges:          args.RepoStatus.OpenChanges,
			InitialBranchesSnapshot: branchesSnapshot,
			InitialConfigSnapshot:   args.Repo.ConfigSnapshot,
			InitialStashSize:        stashSize,
			Lineage:                 args.Config.FullConfig.Lineage,
			PushHook:                args.Config.FullConfig.PushHook,
			RootDir:                 args.Repo.RootDir,
			Run:                     args.Repo.Runner,
			Verbose:                 args.Verbose,
		})
		if err != nil || exit {
			return branchesSnapshot, stashSize, exit, err
		}
	}
	if args.ValidateNoOpenChanges {
		err = validate.NoOpenChanges(args.RepoStatus.OpenChanges)
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
	}
	if args.Fetch {
		var remotes gitdomain.Remotes
		remotes, err := args.Repo.Runner.Backend.Remotes()
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline.Bool() {
			err = args.Repo.Runner.Frontend.Fetch()
			if err != nil {
				return branchesSnapshot, stashSize, false, err
			}
		}
		// must always reload the snapshot here because we fetched updates from the remote
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
	}
	if branchesSnapshot.IsEmpty() {
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
	}
	if args.ValidateIsConfigured {
		err = validate.IsConfigured(&args.Repo.Runner.Backend, args.Config, branchesSnapshot.Branches.LocalBranches().Names(), &args.DialogTestInputs)
	}
	return branchesSnapshot, stashSize, false, err
}

type LoadRepoSnapshotArgs struct {
	Config                *config.Config
	DialogTestInputs      components.TestInputs
	Fetch                 bool
	HandleUnfinishedState bool
	Repo                  *OpenRepoResult
	RepoStatus            gitdomain.RepoStatus
	ValidateIsConfigured  bool
	ValidateNoOpenChanges bool
	Verbose               bool
}
