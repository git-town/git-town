package interpreter

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/print"
	"github.com/git-town/git-town/v17/internal/messages"
	"github.com/git-town/git-town/v17/internal/undo"
	lightInterpreter "github.com/git-town/git-town/v17/internal/vm/interpreter/light"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// autoUndo performs an automatic undo of the current Git Town command.
//
// Some Git Town opcodes can indicate that they auto-undo the entire Git Town command that they are a part of
// should they fail.
func autoUndo(opcode shared.Opcode, runErr error, args ExecuteArgs) error {
	print.Error(fmt.Errorf(messages.RunAutoUndo, runErr.Error()))
	undoProgram, err := undo.CreateUndoForRunningProgram(undo.CreateUndoProgramArgs{
		Backend:        args.Backend,
		Config:         args.Config,
		DryRun:         args.Config.NormalConfig.DryRun,
		Git:            args.Git,
		HasOpenChanges: false,
		NoPushHook:     args.Config.NormalConfig.NoPushHook(),
		RunState:       args.RunState,
	})
	if err != nil {
		return err
	}
	lightInterpreter.Execute(lightInterpreter.ExecuteArgs{
		Backend:       args.Backend,
		Config:        args.Config,
		Connector:     args.Connector,
		FinalMessages: args.FinalMessages,
		Frontend:      args.Frontend,
		Git:           args.Git,
		Prog:          undoProgram,
	})
	return opcode.AutomaticUndoError()
}
