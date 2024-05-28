package execute

import (
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/validate"
)

// LoadRepoSnapshot loads the initial snapshot of the Git repo.
func LoadRepoSnapshot(args LoadRepoSnapshotArgs) (gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	// 1. load saved runstate
	// 2. has saved runstate: validate config (will always work, okay to ask user), continue runstate, exit
	// 3. no saved runstate: do the normal thing - fetch, snapshot, validate config, execute business logic

	// handle unfinished state
	if args.HandleUnfinishedState {
		exit, err := validate.HandleUnfinishedState(validate.UnfinishedStateArgs{
			Backend:           args.Repo.Backend,
			CommandsCounter:   args.Repo.CommandsCounter,
			Connector:         nil,
			DialogTestInputs:  args.DialogTestInputs,
			FinalMessages:     args.Repo.FinalMessages,
			Frontend:          args.Repo.Frontend,
			HasOpenChanges:    args.RepoStatus.OpenChanges,
			InitialBranch:     args.InitialBranch,
			PushHook:          args.UnvalidatedConfig.Config.PushHook,
			RepoStatus:        args.RepoStatus,
			RootDir:           args.Repo.RootDir,
			UnvalidatedConfig: args.UnvalidatedConfig,
			Verbose:           args.Verbose,
		})
		if err != nil || exit {
			return gitdomain.EmptyBranchesSnapshot(), 0, exit, err
		}
	}
	var err error
	if args.ValidateNoOpenChanges {
		err = validate.NoOpenChanges(args.RepoStatus.OpenChanges)
		if err != nil {
			return gitdomain.EmptyBranchesSnapshot(), 0, false, err
		}
	}
	if args.Fetch {
		var remotes gitdomain.Remotes
		remotes, err := args.Backend.Remotes()
		if err != nil {
			return gitdomain.EmptyBranchesSnapshot(), 0, false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline.Bool() {
			err = args.Frontend.Fetch()
			if err != nil {
				return gitdomain.EmptyBranchesSnapshot(), 0, false, err
			}
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
	return branchesSnapshot, stashSize, false, err
}

type LoadRepoSnapshotArgs struct {
	Backend               git.BackendCommands
	CommandsCounter       gohacks.Counter
	ConfigSnapshot        undoconfig.ConfigSnapshot
	DialogTestInputs      components.TestInputs
	Fetch                 bool
	FinalMessages         stringslice.Collector
	Frontend              git.FrontendCommands
	HandleUnfinishedState bool
	InitialBranch         gitdomain.LocalBranchName
	Repo                  OpenRepoResult
	RepoStatus            gitdomain.RepoStatus
	RootDir               gitdomain.RepoRootDir
	UnvalidatedConfig     config.UnvalidatedConfig
	ValidateNoOpenChanges bool
	Verbose               bool
}
