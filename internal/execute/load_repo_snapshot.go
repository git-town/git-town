package execute

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	"github.com/git-town/git-town/v21/internal/validate"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// LoadRepoSnapshot loads the initial snapshot of the Git repo.
func LoadRepoSnapshot(args LoadRepoSnapshotArgs) (gitdomain.BranchesSnapshot, gitdomain.StashSize, Option[gitdomain.BranchInfos], dialogdomain.Exit, error) {
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
			Connector:         args.Connector,
			Detached:          args.Detached,
			FinalMessages:     args.Repo.FinalMessages,
			Frontend:          args.Repo.Frontend,
			Git:               args.Git,
			HasOpenChanges:    args.RepoStatus.OpenChanges,
			Inputs:            args.Inputs,
			PushHook:          args.UnvalidatedConfig.NormalConfig.PushHook,
			RepoStatus:        args.RepoStatus,
			RootDir:           args.Repo.RootDir,
			RunState:          runStateOpt,
			UnvalidatedConfig: args.UnvalidatedConfig,
		})
		if err != nil || exit {
			return gitdomain.EmptyBranchesSnapshot(), 0, previousBranchInfos, exit, err
		}
	}
	if args.ValidateNoOpenChanges {
		if err = validate.NoOpenChanges(args.RepoStatus.OpenChanges); err != nil {
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
			if err = args.Git.Fetch(args.Frontend, args.UnvalidatedConfig.NormalConfig.SyncTags); err != nil {
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
	Backend               subshelldomain.RunnerQuerier
	CommandsCounter       Mutable[gohacks.Counter]
	ConfigSnapshot        undoconfig.ConfigSnapshot
	Connector             Option[forgedomain.Connector]
	Detached              configdomain.Detached
	Fetch                 bool
	FinalMessages         stringslice.Collector
	Frontend              subshelldomain.Runner
	Git                   git.Commands
	HandleUnfinishedState bool
	Inputs                dialogcomponents.Inputs
	Repo                  OpenRepoResult
	RepoStatus            gitdomain.RepoStatus
	RootDir               gitdomain.RepoRootDir
	UnvalidatedConfig     config.UnvalidatedConfig
	ValidateNoOpenChanges bool
}
