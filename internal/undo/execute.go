package undo

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/print"
	"github.com/git-town/git-town/v17/internal/config"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/gohacks"
	"github.com/git-town/git-town/v17/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v17/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v17/internal/messages"
	lightInterpreter "github.com/git-town/git-town/v17/internal/vm/interpreter/light"
	"github.com/git-town/git-town/v17/internal/vm/runstate"
	"github.com/git-town/git-town/v17/internal/vm/statefile"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

// undoes the persisted runstate
func Execute(args ExecuteArgs) error {
	if args.RunState.DryRun {
		return nil
	}
	program := CreateUndoForFinishedProgram(CreateUndoProgramArgs{
		Backend:        args.Backend,
		Config:         args.Config,
		DryRun:         args.RunState.DryRun,
		Git:            args.Git,
		HasOpenChanges: args.HasOpenChanges,
		NoPushHook:     args.Config.NormalConfig.NoPushHook(),
		RunState:       args.RunState,
	})
	lightInterpreter.Execute(lightInterpreter.ExecuteArgs{
		Backend:       args.Backend,
		Config:        args.Config,
		Connector:     args.Connector,
		FinalMessages: args.FinalMessages,
		Frontend:      args.Frontend,
		Git:           args.Git,
		Prog:          program,
	})
	err := statefile.Delete(args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateDeleteProblem, err)
	}
	print.Footer(args.Verbose, args.CommandsCounter.Immutable(), args.FinalMessages.Result())
	return nil
}

type ExecuteArgs struct {
	Backend          gitdomain.RunnerQuerier
	CommandsCounter  Mutable[gohacks.Counter]
	Config           config.ValidatedConfig
	Connector        Option[hostingdomain.Connector]
	FinalMessages    stringslice.Collector
	Frontend         gitdomain.Runner
	Git              git.Commands
	HasOpenChanges   bool
	InitialStashSize gitdomain.StashSize
	RootDir          gitdomain.RepoRootDir
	RunState         runstate.RunState
	Verbose          configdomain.Verbose
}
