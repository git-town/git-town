package undo

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/print"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	lightInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/light"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/git-town/git-town/v12/src/vm/statefile"
)

// undoes the persisted runstate
func Execute(args ExecuteArgs) error {
	fmt.Println("111111111111111111", args.RunState)
	if args.RunState.DryRun {
		return nil
	}
	program := CreateUndoForFinishedProgram(CreateUndoProgramArgs{
		DryRun:         args.Runner.Config.DryRun,
		HasOpenChanges: args.HasOpenChanges,
		NoPushHook:     args.FullConfig.NoPushHook(),
		Run:            args.Runner,
		RunState:       args.RunState,
	})
	fmt.Println("22222222222222222222", program)
	lightInterpreter.Execute(program, args.Runner, args.Lineage)
	err := statefile.Delete(args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateDeleteProblem, err)
	}
	print.Footer(args.Verbose, args.Runner.CommandsCounter.Count(), args.Runner.FinalMessages.Result())
	return nil
}

type ExecuteArgs struct {
	FullConfig       *configdomain.FullConfig
	HasOpenChanges   bool
	InitialStashSize gitdomain.StashSize
	Lineage          configdomain.Lineage
	RootDir          gitdomain.RepoRootDir
	RunState         runstate.RunState
	Runner           *git.ProdRunner
	Verbose          bool
}
