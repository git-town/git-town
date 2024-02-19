package interpreter

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/print"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

// autoUndo performs an automatic undo of the current Git Town command.
//
// Some Git Town opcodes can indicate that they auto-undo the entire Git Town command that they are a part of
// should they fail.
func autoUndo(opcode shared.Opcode, runErr error, args ExecuteArgs) error {
	print.Error(fmt.Errorf(messages.RunAutoUndo, runErr.Error()))
	err := undo.Execute(undo.ExecuteArgs{
		FullConfig:       args.FullConfig,
		HasOpenChanges:   args.HasOpenChanges,
		InitialStashSize: args.InitialStashSize,
		Lineage:          args.Lineage,
		RootDir:          args.RootDir,
		RunState:         *args.RunState,
		Runner:           args.Run,
		Verbose:          args.Verbose,
	})
	if err != nil {
		return fmt.Errorf(messages.RunstateAbortOpcodeProblem, err)
	}
	return opcode.CreateAutomaticUndoError()
}
