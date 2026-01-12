package configinterpreter

import (
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type FinishedArgs struct {
	Backend               subshelldomain.RunnerQuerier
	BeginBranchesSnapshot Option[gitdomain.BranchesSnapshot]
	BeginConfigSnapshot   configdomain.BeginConfigSnapshot
	Command               string
	CommandsCounter       Mutable[gohacks.Counter]
	ConfigDir             configdomain.ConfigDirRepo
	FinalMessages         stringslice.Collector
	Git                   git.Commands
	RootDir               gitdomain.RepoRootDir
	TouchedBranches       []gitdomain.BranchName
	Verbose               configdomain.Verbose
}

// Finished is called when a Git Town command that only changes configuration has finished successfully.
func Finished(args FinishedArgs) error {
	var endBranchesSnapshot Option[gitdomain.BranchesSnapshot]
	if args.BeginBranchesSnapshot.IsSome() {
		snapshot, err := args.Git.BranchesSnapshot(args.Backend)
		if err != nil {
			return err
		}
		endBranchesSnapshot = Some(snapshot)
	}
	globalSnapshot, err := gitconfig.LoadSnapshot(args.Backend, Some(configdomain.ConfigScopeGlobal), configdomain.UpdateOutdatedNo)
	if err != nil {
		return err
	}
	localSnapshot, err := gitconfig.LoadSnapshot(args.Backend, Some(configdomain.ConfigScopeLocal), configdomain.UpdateOutdatedNo)
	if err != nil {
		return err
	}
	configSnapshot := configdomain.EndConfigSnapshot{
		Global: globalSnapshot,
		Local:  localSnapshot,
	}
	runState := runstate.RunState{
		AbortProgram:             program.Program{},
		BeginBranchesSnapshot:    args.BeginBranchesSnapshot.GetOrZero(),
		BeginConfigSnapshot:      args.BeginConfigSnapshot,
		BeginStashSize:           0,
		Command:                  args.Command,
		DryRun:                   false,
		EndBranchesSnapshot:      endBranchesSnapshot,
		EndConfigSnapshot:        Some(configSnapshot),
		EndStashSize:             None[gitdomain.StashSize](),
		FinalUndoProgram:         program.Program{},
		BranchInfosLastRun:       None[gitdomain.BranchInfos](),
		RunProgram:               program.Program{},
		TouchedBranches:          args.TouchedBranches,
		UndoablePerennialCommits: gitdomain.SHAs{},
		UndoAPIProgram:           program.Program{},
		UnfinishedDetails:        MutableNone[runstate.UnfinishedRunStateDetails](),
	}
	print.Footer(args.Verbose, args.CommandsCounter.Immutable(), args.FinalMessages.Result())
	return runstate.Save(runState, args.RunstatePath)
}
