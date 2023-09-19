package runvm

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/persistence"
	"github.com/git-town/git-town/v9/src/steps"
)

// errored is called when the given step has resulted in the given error.
func errored(step steps.Step, runErr error, args ExecuteArgs) error {
	args.RunState.AbortStepList.Append(step.CreateAbortSteps()...)
	if step.ShouldAutomaticallyAbortOnError() {
		return autoAbort(step, runErr, args)
	}
	args.RunState.RunStepList.Prepend(step.CreateContinueSteps()...)
	err := args.RunState.MarkAsUnfinished(&args.Run.Backend)
	if err != nil {
		return err
	}
	currentBranch, err := args.Run.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	rebasing, err := args.Run.Backend.HasRebaseInProgress()
	if err != nil {
		return err
	}
	if args.RunState.Command == "sync" && !(rebasing && args.Run.Config.IsMainBranch(currentBranch)) {
		args.RunState.UnfinishedDetails.CanSkip = true
	}
	err = persistence.Save(args.RunState, args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	message := runErr.Error() + messages.AbortContinueGuidance
	if args.RunState.UnfinishedDetails.CanSkip {
		message += messages.ContinueSkipGuidance
	}
	message += "\n"
	return fmt.Errorf(message)
}
