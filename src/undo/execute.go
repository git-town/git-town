package undo

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/messages"
	lightInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/light"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/git-town/git-town/v14/src/vm/statefile"
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
	fmt.Println("444444444444444444", program)
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
	print.Footer(args.Verbose, args.CommandsCounter.Count(), args.FinalMessages.Result())
	return nil
}

type ExecuteArgs struct {
	Backend          gitdomain.RunnerQuerier
	CommandsCounter  gohacks.Counter
	Config           config.ValidatedConfig
	FinalMessages    stringslice.Collector
	Frontend         gitdomain.Runner
	Git              git.Commands
	HasOpenChanges   bool
	InitialStashSize gitdomain.StashSize
	RootDir          gitdomain.RepoRootDir
	RunState         runstate.RunState
	Verbose          bool
}
