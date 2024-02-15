package execute

import (
	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/validate"
)

// LoadRepoSnapshot loads the initial snapshot of the Git repo.
func LoadRepoSnapshot(args LoadRepoSnapshotArgs) (gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	var branchesSnapshot gitdomain.BranchesSnapshot
	var err error
	stashSize, err := args.Repo.Runner.Backend.StashSize()
	if err != nil {
		return gitdomain.EmptyBranchesSnapshot(), stashSize, false, err
	}
	if args.HandleUnfinishedState {
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
		exit, err := validate.HandleUnfinishedState(validate.UnfinishedStateArgs{
			Connector:             nil,
			DialogTestInputs:      args.DialogTestInputs,
			Verbose:               args.Verbose,
			InitialConfigSnapshot: args.Repo.ConfigSnapshot,
			InitialStashSize:      stashSize,
			Lineage:               args.Lineage,
			PushHook:              args.PushHook,
			RootDir:               args.Repo.RootDir,
			Run:                   args.Repo.Runner,
		})
		if err != nil || exit {
			return branchesSnapshot, stashSize, exit, err
		}
	}
	if args.ValidateNoOpenChanges {
		repoStatus, err := args.Repo.Runner.Backend.RepoStatus()
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
		err = validate.NoOpenChanges(repoStatus.OpenChanges)
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
		err = validate.IsConfigured(&args.Repo.Runner.Backend, args.FullConfig, branchesSnapshot.Branches.LocalBranches().Names(), &args.DialogTestInputs)
	}
	return branchesSnapshot, stashSize, false, err
}

type LoadRepoSnapshotArgs struct {
	*configdomain.FullConfig
	DialogTestInputs      components.TestInputs
	Fetch                 bool
	HandleUnfinishedState bool
	Repo                  *OpenRepoResult
	ValidateIsConfigured  bool
	ValidateNoOpenChanges bool
	Verbose               bool
}
