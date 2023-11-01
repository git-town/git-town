package interpreter

import (
	"fmt"

	"github.com/git-town/git-town/v10/src/cli/print"
	"github.com/git-town/git-town/v10/src/messages"
	"github.com/git-town/git-town/v10/src/vm/shared"
)

// autoAbort performs an automatic abort of the current Git Town command.
//
// Some Git Town opcodes can indicate that they auto-abort the entire Git Town command that they are a part of
// should they fail.
func autoAbort(opcode shared.Opcode, runErr error, args ExecuteArgs) error {
	print.Error(fmt.Errorf(messages.RunAutoAborting, runErr.Error()))
	abortRunState := args.RunState.CreateAbortRunState()
	err := Execute(ExecuteArgs{
		RunState:                &abortRunState,
		Run:                     args.Run,
		Connector:               args.Connector,
		Verbose:                 args.Verbose,
		RootDir:                 args.RootDir,
		Lineage:                 args.Lineage,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
		InitialStashSnapshot:    args.InitialStashSnapshot,
		NoPushHook:              args.NoPushHook,
	})
	if err != nil {
		return fmt.Errorf(messages.RunstateAbortOpcodeProblem, err)
	}
	return opcode.CreateAutomaticAbortError()
}
