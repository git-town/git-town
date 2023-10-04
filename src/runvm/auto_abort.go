package runvm

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/step"
)

// autoAbort performs an automatic abort of the current Git Town command.
//
// Some Git Town steps can indicate that they auto-abort the entire Git Town command that they are a part of
// should they fail.
func autoAbort(step step.Step, runErr error, args ExecuteArgs) error {
	cli.PrintError(fmt.Errorf(messages.RunAutoAborting, runErr.Error()))
	abortRunState := args.RunState.CreateAbortRunState()
	err := Execute(ExecuteArgs{
		RunState:                &abortRunState,
		Run:                     args.Run,
		Connector:               args.Connector,
		RootDir:                 args.RootDir,
		Lineage:                 args.Lineage,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
		InitialStashSnapshot:    args.InitialStashSnapshot,
		NoPushHook:              args.NoPushHook,
	})
	if err != nil {
		return fmt.Errorf(messages.RunstateAbortStepProblem, err)
	}
	return step.CreateAutomaticAbortError()
}
