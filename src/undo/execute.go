package undo

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/config/configdomain"
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
		Config:         args.Config,
		DryRun:         args.Config.DryRun,
		HasOpenChanges: args.HasOpenChanges,
		NoPushHook:     args.Config.Config.NoPushHook(),
		RunState:       args.RunState,
	})
	lightInterpreter.Execute(lightInterpreter.ExecuteArgs{
		Backend:       args.Backend,
		Config:        args.Config,
		FinalMessages: args.FinalMessages,
		Frontend:      args.Frontend,
		Lineage:       args.Lineage,
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
	Backend          git.BackendCommands
	CommandsCounter  *gohacks.Counter
	Config           configdomain.ValidatedConfig
	FinalMessages    *stringslice.Collector
	Frontend         git.FrontendCommands
	HasOpenChanges   bool
	InitialStashSize gitdomain.StashSize
	Lineage          configdomain.Lineage
	RootDir          gitdomain.RepoRootDir
	RunState         runstate.RunState
	Verbose          bool
}
