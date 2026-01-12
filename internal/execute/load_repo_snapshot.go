package execute

import (
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v22/internal/validate"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// LoadRepoSnapshot loads the initial snapshot of the Git repo.
func LoadRepoSnapshot(args LoadRepoSnapshotArgs) (gitdomain.BranchesSnapshot, gitdomain.StashSize, Option[gitdomain.BranchInfos], configdomain.ProgramFlow, error) {
	runstatePath := runstate.NewRunstatePath(args.Repo.ConfigDir)
	runStateOpt, err := runstate.Load(runstatePath)
	if err != nil {
		return gitdomain.EmptyBranchesSnapshot(), 0, None[gitdomain.BranchInfos](), configdomain.ProgramFlowExit, err
	}
	previousBranchInfos := None[gitdomain.BranchInfos]()
	if runstate, hasRunstate := runStateOpt.Get(); hasRunstate {
		if endSnapshot, hasEndSnapshot := runstate.EndBranchesSnapshot.Get(); hasEndSnapshot {
			runstate.BranchInfosLastRun = Some(endSnapshot.Branches)
		}
		previousBranchInfos = runstate.BranchInfosLastRun
	}
	if args.HandleUnfinishedState {
		flow, err := validate.HandleUnfinishedState(validate.UnfinishedStateArgs{
			Backend:           args.Repo.Backend,
			CommandsCounter:   args.Repo.CommandsCounter,
			ConfigDir:         args.Repo.ConfigDir,
			Connector:         args.Connector,
			FinalMessages:     args.Repo.FinalMessages,
			Frontend:          args.Repo.Frontend,
			Git:               args.Git,
			HasOpenChanges:    args.RepoStatus.OpenChanges,
			Inputs:            args.Inputs,
			PushHook:          args.UnvalidatedConfig.NormalConfig.PushHook,
			RepoStatus:        args.RepoStatus,
			RunState:          runStateOpt,
			UnvalidatedConfig: args.UnvalidatedConfig,
		})
		if err != nil {
			return gitdomain.EmptyBranchesSnapshot(), 0, previousBranchInfos, configdomain.ProgramFlowExit, err
		}
		switch flow {
		case configdomain.ProgramFlowContinue:
		case configdomain.ProgramFlowExit, configdomain.ProgramFlowRestart:
			return gitdomain.EmptyBranchesSnapshot(), 0, previousBranchInfos, flow, nil
		}
	}
	if args.ValidateNoOpenChanges {
		if err = validate.NoOpenChanges(args.RepoStatus.OpenChanges); err != nil {
			return gitdomain.EmptyBranchesSnapshot(), 0, previousBranchInfos, configdomain.ProgramFlowExit, err
		}
	}
	if args.Fetch {
		var remotes gitdomain.Remotes
		remotes, err := args.Git.Remotes(args.Repo.Backend)
		if err != nil {
			return gitdomain.EmptyBranchesSnapshot(), 0, previousBranchInfos, configdomain.ProgramFlowExit, err
		}
		if remotes.HasRemote(args.UnvalidatedConfig.NormalConfig.DevRemote) && args.Repo.IsOffline.IsOnline() {
			if err = args.Git.Fetch(args.Frontend, args.UnvalidatedConfig.NormalConfig.SyncTags); err != nil {
				return gitdomain.EmptyBranchesSnapshot(), 0, previousBranchInfos, configdomain.ProgramFlowExit, err
			}
		}
	}
	stashSize, err := args.Repo.Git.StashSize(args.Repo.Backend)
	if err != nil {
		return gitdomain.EmptyBranchesSnapshot(), stashSize, previousBranchInfos, configdomain.ProgramFlowExit, err
	}
	branchesSnapshot, err := args.Repo.Git.BranchesSnapshot(args.Repo.Backend)
	return branchesSnapshot, stashSize, previousBranchInfos, configdomain.ProgramFlowContinue, err
}

type LoadRepoSnapshotArgs struct {
	Backend               subshelldomain.RunnerQuerier
	CommandsCounter       Mutable[gohacks.Counter]
	ConfigSnapshot        configdomain.BeginConfigSnapshot
	Connector             Option[forgedomain.Connector]
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
