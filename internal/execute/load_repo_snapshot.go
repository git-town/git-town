package execute

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	"github.com/git-town/git-town/v21/internal/validate"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// LoadRepoSnapshot loads the initial snapshot of the Git repo.
func LoadRepoSnapshot(args LoadRepoSnapshotArgs) (gitdomain.BranchesSnapshot, gitdomain.StashSize, Option[gitdomain.BranchInfos], bool, error) {
	runStateOpt, err := runstate.Load(args.RootDir)
	if err != nil {
		return gitdomain.EmptyBranchesSnapshot(), 0, None[gitdomain.BranchInfos](), false, err
	}
	previousBranchInfos := None[gitdomain.BranchInfos]()
	if runstate, hasRunstate := runStateOpt.Get(); hasRunstate {
		if endSnapshot, hasEndSnapshot := runstate.EndBranchesSnapshot.Get(); hasEndSnapshot {
			runstate.BranchInfosLastRun = Some(endSnapshot.Branches)
		}
		previousBranchInfos = runstate.BranchInfosLastRun
	}
	if args.HandleUnfinishedState {
		exit, err := validate.HandleUnfinishedState(validate.UnfinishedStateArgs{
			Backend:           args.Repo.Backend,
			CommandsCounter:   args.Repo.CommandsCounter,
			Connector:         None[forgedomain.Connector](),
			Detached:          args.Detached,
			DialogTestInputs:  args.DialogTestInputs,
			FinalMessages:     args.Repo.FinalMessages,
			Frontend:          args.Repo.Frontend,
			Git:               args.Git,
			HasOpenChanges:    args.RepoStatus.OpenChanges,
			PushHook:          args.UnvalidatedConfig.NormalConfig.PushHook,
			RepoStatus:        args.RepoStatus,
			RootDir:           args.Repo.RootDir,
			RunState:          runStateOpt,
			UnvalidatedConfig: args.UnvalidatedConfig,
			Verbose:           args.Verbose,
		})
		if err != nil || exit {
			return gitdomain.EmptyBranchesSnapshot(), 0, previousBranchInfos, exit, err
		}
	}
	if args.ValidateNoOpenChanges {
		err = validate.NoOpenChanges(args.RepoStatus.OpenChanges)
		if err != nil {
			return gitdomain.EmptyBranchesSnapshot(), 0, previousBranchInfos, false, err
		}
	}
	if args.Fetch {
		var remotes gitdomain.Remotes
		remotes, err := args.Git.Remotes(args.Repo.Backend)
		if err != nil {
			return gitdomain.EmptyBranchesSnapshot(), 0, previousBranchInfos, false, err
		}
		if remotes.HasRemote(args.UnvalidatedConfig.NormalConfig.DevRemote) && args.Repo.IsOffline.IsOnline() {
			err = args.Git.Fetch(args.Frontend, args.UnvalidatedConfig.NormalConfig.SyncTags)
			if err != nil {
				return gitdomain.EmptyBranchesSnapshot(), 0, previousBranchInfos, false, err
			}
		}
	}
	stashSize, err := args.Repo.Git.StashSize(args.Repo.Backend)
	if err != nil {
		return gitdomain.EmptyBranchesSnapshot(), stashSize, previousBranchInfos, false, err
	}
	branchesSnapshot, err := args.Repo.Git.BranchesSnapshot(args.Repo.Backend)
	return branchesSnapshot, stashSize, previousBranchInfos, false, err
}

type LoadRepoSnapshotArgs struct {
	Backend               gitdomain.RunnerQuerier
	CommandsCounter       Mutable[gohacks.Counter]
	ConfigSnapshot        undoconfig.ConfigSnapshot
	Detached              configdomain.Detached
	DialogTestInputs      components.TestInputs
	Fetch                 bool
	FinalMessages         stringslice.Collector
	Frontend              gitdomain.Runner
	Git                   git.Commands
	HandleUnfinishedState bool
	Repo                  OpenRepoResult
	RepoStatus            gitdomain.RepoStatus
	RootDir               gitdomain.RepoRootDir
	UnvalidatedConfig     config.UnvalidatedConfig
	ValidateNoOpenChanges bool
	Verbose               configdomain.Verbose
}
