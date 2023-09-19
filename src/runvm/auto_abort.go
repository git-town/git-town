package runvm

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/steps"
)

// autoAbort is called when a step that produced an error triggers an auto-abort.
func autoAbort(step steps.Step, runErr error, args ExecuteArgs) error {
	cli.PrintError(fmt.Errorf(messages.RunAutoAborting, runErr.Error()))
	abortRunState := args.RunState.CreateAbortRunState()
	err := Execute(ExecuteArgs{
		RunState:  &abortRunState,
		Run:       args.Run,
		Connector: args.Connector,
		RootDir:   args.RootDir,
		Lineage:   args.Lineage,
	})
	if err != nil {
		return fmt.Errorf(messages.RunstateAbortStepProblem, err)
	}
	return step.CreateAutomaticAbortError()
}
