package interpreter

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/undo"
	"github.com/git-town/git-town/v9/src/vm/persistence"
)

// finished is called when executing all steps has successfully finished.
func finished(args ExecuteArgs) error {
	args.RunState.MarkAsFinished()
	undoProgram, err := undo.CreateUndoProgram(undo.CreateUndoProgramArgs{
		Run:                      args.Run,
		InitialBranchesSnapshot:  args.InitialBranchesSnapshot,
		InitialConfigSnapshot:    args.InitialConfigSnapshot,
		InitialStashSnapshot:     args.InitialStashSnapshot,
		NoPushHook:               args.NoPushHook,
		UndoablePerennialCommits: args.RunState.UndoablePerennialCommits,
	})
	if err != nil {
		return err
	}
	args.RunState.UndoProgram.AddProgram(undoProgram)
	args.RunState.UndoProgram.AddProgram(args.RunState.FinalUndoProgram)
	if args.RunState.IsAbort || args.RunState.IsUndo {
		err := persistence.Delete(args.RootDir)
		if err != nil {
			return fmt.Errorf(messages.RunstateDeleteProblem, err)
		}
	} else {
		err := persistence.Save(args.RunState, args.RootDir)
		if err != nil {
			return fmt.Errorf(messages.RunstateSaveProblem, err)
		}
	}
	PrintFooter(args.Debug, args.Run.CommandsCounter.Count(), args.Run.FinalMessages.Result())
	return nil
}
