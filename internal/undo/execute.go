package undo

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/cli/print"
	"github.com/git-town/git-town/v15/internal/config"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/gohacks"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v15/internal/messages"
	lightInterpreter "github.com/git-town/git-town/v15/internal/vm/interpreter/light"
	"github.com/git-town/git-town/v15/internal/vm/runstate"
	"github.com/git-town/git-town/v15/internal/vm/statefile"
)

// undoes the persisted runstate
func Execute(args ExecuteArgs) error {
	if args.RunState.DryRun {
		return nil
	}
	program := CreateUndoForFinishedProgram(CreateUndoProgramArgs{
		Backend:        args.Backend,
		Config:         args.Config.Config,
		DryRun:         args.RunState.DryRun,
		Git:            args.Git,
		HasOpenChanges: args.HasOpenChanges,
		NoPushHook:     args.Config.Config.NoPushHook(),
		RunState:       args.RunState,
	})
	lightInterpreter.Execute(lightInterpreter.ExecuteArgs{
		Backend:       args.Backend,
		Config:        args.Config,
		FinalMessages: args.FinalMessages,
		Frontend:      args.Frontend,
		Git:           args.Git,
		Prog:          program,
	})
	err := statefile.Delete(args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateDeleteProblem, err)
	}
	print.Footer(args.Verbose, args.CommandsCounter.Get(), args.FinalMessages.Result())
	return nil
}

type ExecuteArgs struct {
	Backend          gitdomain.RunnerQuerier
	CommandsCounter  Mutable[gohacks.Counter]
	Config           config.ValidatedConfig
	FinalMessages    stringslice.Collector
	Frontend         gitdomain.Runner
	Git              git.Commands
	HasOpenChanges   bool
	InitialStashSize gitdomain.StashSize
	RootDir          gitdomain.RepoRootDir
	RunState         runstate.RunState
	Verbose          configdomain.Verbose
}
