package execute

import (
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/validate"
)

// LoadRepoSnapshot loads the initial snapshot of the Git repo.
func LoadRepoSnapshot(args LoadRepoSnapshotArgs) (gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	var err error
	if args.HandleUnfinishedState {
		exit, err := validate.HandleUnfinishedState(validate.UnfinishedStateArgs{
			Backend:          args.Repo.Backend,
			CommandsCounter:  args.Repo.CommandsCounter,
			Config:           args.Config,
			Connector:        nil,
			DialogTestInputs: args.DialogTestInputs,
			FinalMessages:    args.Repo.FinalMessages,
			Frontend:         args.Repo.Frontend,
			HasOpenChanges:   args.RepoStatus.OpenChanges,
			Lineage:          args.Config.Config.Lineage,
			PushHook:         args.Config.Config.PushHook,
			RootDir:          args.Repo.RootDir,
			Verbose:          args.Verbose,
		})
		if err != nil || exit {
			return gitdomain.EmptyBranchesSnapshot(), 0, exit, err
		}
	}
	stashSize, err := args.Repo.Backend.StashSize()
	if err != nil {
		return gitdomain.EmptyBranchesSnapshot(), stashSize, false, err
	}
	branchesSnapshot, err := args.Repo.Backend.BranchesSnapshot()
	if err != nil {
		return branchesSnapshot, stashSize, false, err
	}
	if args.ValidateNoOpenChanges {
		err = validate.NoOpenChanges(args.RepoStatus.OpenChanges)
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
	}
	if args.Fetch {
		var remotes gitdomain.Remotes
		remotes, err := args.Repo.Backend.Remotes()
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline.Bool() {
			err = args.Repo.Frontend.Fetch()
			if err != nil {
				return branchesSnapshot, stashSize, false, err
			}
		}
		// must always reload the snapshot here because we fetched updates from the remote
		branchesSnapshot, err = args.Repo.Backend.BranchesSnapshot()
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
	}
	if branchesSnapshot.IsEmpty() {
		branchesSnapshot, err = args.Repo.Backend.BranchesSnapshot()
		if err != nil {
			return branchesSnapshot, stashSize, false, err
		}
	}
	return branchesSnapshot, stashSize, false, err
}

type LoadRepoSnapshotArgs struct {
	Config                config.Config
	DialogTestInputs      components.TestInputs
	Fetch                 bool
	HandleUnfinishedState bool
	Repo                  OpenRepoResult
	RepoStatus            gitdomain.RepoStatus
	ValidateNoOpenChanges bool
	Verbose               bool
}
