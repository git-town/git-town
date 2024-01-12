package execute

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/validate"
)

// LoadRepoSnapshot loads the initial snapshot of the Git repo.
func LoadRepoSnapshot(args LoadBranchesArgs) (gitdomain.BranchesStatus, gitdomain.StashSize, dialog.TestInputs, bool, error) {
	var branchesSnapshot gitdomain.BranchesStatus
	dialogInputs := dialog.LoadTestInputs(os.Environ())
	var err error
	stashSnapshot, err := args.Repo.Runner.Backend.StashSize()
	if err != nil {
		return gitdomain.EmptyBranchesSnapshot(), stashSnapshot, dialogInputs, false, err
	}
	if args.HandleUnfinishedState {
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return branchesSnapshot, stashSnapshot, dialogInputs, false, err
		}
		exit, err := validate.HandleUnfinishedState(validate.UnfinishedStateArgs{
			Connector:               nil,
			Verbose:                 args.Verbose,
			InitialBranchesSnapshot: branchesSnapshot,
			InitialConfigSnapshot:   args.Repo.ConfigSnapshot,
			InitialStashSnapshot:    stashSnapshot,
			Lineage:                 args.Lineage,
			PushHook:                args.PushHook,
			RootDir:                 args.Repo.RootDir,
			Run:                     args.Repo.Runner,
		})
		if err != nil || exit {
			return branchesSnapshot, stashSnapshot, dialogInputs, exit, err
		}
	}
	if args.ValidateNoOpenChanges {
		repoStatus, err := args.Repo.Runner.Backend.RepoStatus()
		if err != nil {
			return branchesSnapshot, stashSnapshot, dialogInputs, false, err
		}
		err = validate.NoOpenChanges(repoStatus.OpenChanges)
		if err != nil {
			return branchesSnapshot, stashSnapshot, dialogInputs, false, err
		}
	}
	if args.Fetch {
		var remotes gitdomain.Remotes
		remotes, err := args.Repo.Runner.Backend.Remotes()
		if err != nil {
			return branchesSnapshot, stashSnapshot, dialogInputs, false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline.Bool() {
			err = args.Repo.Runner.Frontend.Fetch()
			if err != nil {
				return branchesSnapshot, stashSnapshot, dialogInputs, false, err
			}
		}
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return branchesSnapshot, stashSnapshot, dialogInputs, false, err
		}
	}
	if branchesSnapshot.IsEmpty() {
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return branchesSnapshot, stashSnapshot, dialogInputs, false, err
		}
	}
	if args.ValidateIsConfigured {
		err = validate.IsConfigured(&args.Repo.Runner.Backend, args.FullConfig, branchesSnapshot.Branches.LocalBranches().Names(), dialogInputs)
	}
	return branchesSnapshot, stashSnapshot, dialogInputs, false, err
}

type LoadBranchesArgs struct {
	*configdomain.FullConfig
	Repo                  *OpenRepoResult
	Fetch                 bool
	Verbose               bool
	HandleUnfinishedState bool
	ValidateIsConfigured  bool
	ValidateNoOpenChanges bool
}
