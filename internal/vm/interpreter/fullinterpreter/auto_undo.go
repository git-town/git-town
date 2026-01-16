package fullinterpreter

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/undo"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/lightinterpreter"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// autoUndo performs an automatic undo of the current Git Town command.
//
// Some Git Town opcodes can indicate that they auto-undo the entire Git Town command that they are a part of
// should they fail.
func autoUndo(opcode shared.AutoUndoable, runErr error, args ExecuteArgs) error {
	print.Error(fmt.Errorf(messages.RunAutoUndo, runErr.Error()))
	undoProgram, err := undo.CreateUndoForRunningProgram(undo.CreateUndoProgramArgs{
		Backend:        args.Backend,
		Config:         args.Config,
		DryRun:         args.Config.NormalConfig.DryRun,
		FinalMessages:  args.FinalMessages,
		Git:            args.Git,
		HasOpenChanges: false,
		PushHook:       args.Config.NormalConfig.PushHook,
		RunState:       args.RunState,
	})
	if err != nil {
		return err
	}
	lightinterpreter.Execute(lightinterpreter.ExecuteArgs{
		Backend:       args.Backend,
		BranchInfos:   args.RunState.BeginBranchesSnapshot.Branches,
		Config:        args.Config,
		Connector:     args.Connector,
		FinalMessages: args.FinalMessages,
		Frontend:      args.Frontend,
		Git:           args.Git,
		Prog:          undoProgram,
	})
	return opcode.AutomaticUndoError()
}
