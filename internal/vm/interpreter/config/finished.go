package config

import (
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/config/gitconfig"
	"github.com/git-town/git-town/v16/internal/git"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v16/internal/undo/undoconfig"
	"github.com/git-town/git-town/v16/internal/vm/program"
	"github.com/git-town/git-town/v16/internal/vm/runstate"
	"github.com/git-town/git-town/v16/internal/vm/statefile"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

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
	configGitAccess := gitconfig.Access{Runner: args.Backend}
	globalSnapshot, _, err := configGitAccess.LoadGlobal(false)
	if err != nil {
		return err
	}
	localSnapshot, _, err := configGitAccess.LoadLocal(false)
	if err != nil {
		return err
	}
	configSnapshot := undoconfig.ConfigSnapshot{
		Global: globalSnapshot,
		Local:  localSnapshot,
	}
	runState := runstate.RunState{
		AbortProgram:             program.Program{},
		BeginBranchesSnapshot:    args.BeginBranchesSnapshot.GetOrDefault(),
		BeginConfigSnapshot:      args.BeginConfigSnapshot,
		BeginStashSize:           0,
		Command:                  args.Command,
		DryRun:                   false,
		EndBranchesSnapshot:      endBranchesSnapshot,
		EndConfigSnapshot:        Some(configSnapshot),
		EndStashSize:             None[gitdomain.StashSize](),
		FinalUndoProgram:         program.Program{},
		RunProgram:               program.Program{},
		TouchedBranches:          args.TouchedBranches,
		UndoablePerennialCommits: gitdomain.SHAs{},
		UndoAPIProgram:           program.Program{},
		UnfinishedDetails:        MutableNone[runstate.UnfinishedRunStateDetails](),
	}
	print.Footer(args.Verbose, args.CommandsCounter.Get(), args.FinalMessages.Result())
	return statefile.Save(runState, args.RootDir)
}

type FinishedArgs struct {
	Backend               gitdomain.RunnerQuerier
	BeginBranchesSnapshot Option[gitdomain.BranchesSnapshot]
	BeginConfigSnapshot   undoconfig.ConfigSnapshot
	Command               string
	CommandsCounter       Mutable[gohacks.Counter]
	FinalMessages         stringslice.Collector
	Git                   git.Commands
	RootDir               gitdomain.RepoRootDir
	TouchedBranches       []gitdomain.BranchName
	Verbose               configdomain.Verbose
}
