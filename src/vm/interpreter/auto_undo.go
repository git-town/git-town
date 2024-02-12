package interpreter

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/print"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

// autoUndo performs an automatic undo of the current Git Town command.
//
// Some Git Town opcodes can indicate that they auto-undo the entire Git Town command that they are a part of
// should they fail.
func autoUndo(opcode shared.Opcode, runErr error, args ExecuteArgs) error {
	print.Error(fmt.Errorf(messages.RunAutoUndo, runErr.Error()))
	abortRunState := args.RunState.CreateAbortRunState()
	err := Execute(ExecuteArgs{
		Connector:               args.Connector,
		DialogTestInputs:        args.DialogTestInputs,
		FullConfig:              args.FullConfig,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
		InitialStashSize:        args.InitialStashSize,
		RootDir:                 args.RootDir,
		Run:                     args.Run,
		RunState:                &abortRunState,
		Verbose:                 args.Verbose,
	})
	if err != nil {
		return fmt.Errorf(messages.RunstateAbortOpcodeProblem, err)
	}
	return opcode.CreateAutomaticUndoError()
}
