package undo

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
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
		DryRun:         args.Runner.Config.DryRun,
		HasOpenChanges: args.HasOpenChanges,
		NoPushHook:     args.Config.NoPushHook(),
		Run:            args.Runner,
		RunState:       args.RunState,
	})
	lightInterpreter.Execute(program, args.Runner, args.Lineage)
	err := statefile.Delete(args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateDeleteProblem, err)
	}
	print.Footer(args.Verbose, args.Runner.CommandsCounter.Count(), args.Runner.FinalMessages.Result())
	return nil
}

type ExecuteArgs struct {
	Config           configdomain.FullConfig
	HasOpenChanges   bool
	InitialStashSize gitdomain.StashSize
	Lineage          configdomain.Lineage
	RootDir          gitdomain.RepoRootDir
	RunState         runstate.RunState
	Runner           *git.ProdRunner
	Verbose          bool
}
