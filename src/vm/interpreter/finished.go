package interpreter

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/print"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/git-town/git-town/v11/src/undo"
	"github.com/git-town/git-town/v11/src/vm/statefile"
)

// finished is called when executing all steps has successfully finished.
func finished(args ExecuteArgs) error {
	if args.RunState.IsUndo {
		return finishedUndoCommand(args)
	}
	if args.RunState.DryRun {
		return finishedDryRunCommand(args)
	}
	args.RunState.MarkAsFinished()
	undoProgram, err := undo.CreateUndoProgram(undo.CreateUndoProgramArgs{
		DryRun:                   args.RunState.DryRun,
		Run:                      args.Run,
		InitialBranchesSnapshot:  args.InitialBranchesSnapshot,
		InitialConfigSnapshot:    args.InitialConfigSnapshot,
		InitialStashSnapshot:     args.InitialStashSnapshot,
		NoPushHook:               args.NoPushHook(),
		UndoablePerennialCommits: args.RunState.UndoablePerennialCommits,
	})
	if err != nil {
		return err
	}
	args.RunState.UndoProgram.AddProgram(undoProgram)
	args.RunState.UndoProgram.AddProgram(args.RunState.FinalUndoProgram)
	err = statefile.Save(args.RunState, args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	print.Footer(args.Verbose, args.Run.CommandsCounter.Count(), args.Run.FinalMessages.Result())
	return nil
}

func finishedDryRunCommand(args ExecuteArgs) error {
	args.RunState.MarkAsFinished()
	err := statefile.Save(args.RunState, args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	print.Footer(args.Verbose, args.Run.CommandsCounter.Count(), args.Run.FinalMessages.Result())
	return nil
}

func finishedUndoCommand(args ExecuteArgs) error {
	err := statefile.Delete(args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateDeleteProblem, err)
	}
	print.Footer(args.Verbose, args.Run.CommandsCounter.Count(), args.Run.FinalMessages.Result())
	return nil
}
