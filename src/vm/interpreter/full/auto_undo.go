package interpreter

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/print"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo"
	"github.com/git-town/git-town/v12/src/vm/interpreter/light"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

// autoUndo performs an automatic undo of the current Git Town command.
//
// Some Git Town opcodes can indicate that they auto-undo the entire Git Town command that they are a part of
// should they fail.
func autoUndo(opcode shared.Opcode, runErr error, args ExecuteArgs) error {
	print.Error(fmt.Errorf(messages.RunAutoUndo, runErr.Error()))
	undoProgram, err := undo.CreateUndoForRunningProgram(undo.CreateUndoProgramArgs{
		DryRun:         args.Run.DryRun,
		HasOpenChanges: false,
		NoPushHook:     args.FullConfig.NoPushHook(),
		Run:            args.Run,
		RunState:       *args.RunState,
	})
	if err != nil {
		return err
	}
	light.Execute(undoProgram, args.Run, args.Lineage)
	return opcode.CreateAutomaticUndoError()
}
