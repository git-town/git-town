package interpreter

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo"
	lightInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/light"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// autoUndo performs an automatic undo of the current Git Town command.
//
// Some Git Town opcodes can indicate that they auto-undo the entire Git Town command that they are a part of
// should they fail.
func autoUndo(opcode shared.Opcode, runErr error, args ExecuteArgs) error {
	print.Error(fmt.Errorf(messages.RunAutoUndo, runErr.Error()))
	undoProgram, err := undo.CreateUndoForRunningProgram(undo.CreateUndoProgramArgs{
		Backend:        args.Backend,
		Config:         args.Config.Config,
		DryRun:         args.Config.DryRun,
		Git:            args.Git,
		HasOpenChanges: false,
		NoPushHook:     args.Config.Config.NoPushHook(),
		RunState:       args.RunState,
	})
	if err != nil {
		return err
	}
	lightInterpreter.Execute(lightInterpreter.ExecuteArgs{
		Backend:       args.Backend,
		Config:        args.Config,
		FinalMessages: args.FinalMessages,
		Frontend:      args.Frontend,
		Git:           args.Git,
		Prog:          undoProgram,
	})
	return opcode.CreateAutomaticUndoError()
}
