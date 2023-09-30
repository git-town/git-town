package runvm

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/persistence"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/git-town/git-town/v9/src/undo"
)

// errored is called when the given step has resulted in the given error.
func errored(step steps.Step, runErr error, args ExecuteArgs) error {
	args.RunState.AbortSteps.Append(step.CreateAbortSteps()...)
	undoSteps, err := undo.CreateUndoList(undo.CreateUndoListArgs{
		Run:                      args.Run,
		InitialBranchesSnapshot:  args.InitialBranchesSnapshot,
		InitialConfigSnapshot:    args.InitialConfigSnapshot,
		InitialStashSnapshot:     args.InitialStashSnapshot,
		NoPushHook:               args.NoPushHook,
		UndoablePerennialCommits: args.RunState.UndoablePerennialCommits,
	})
	if err != nil {
		return err
	}
	args.RunState.UndoSteps.AppendList(undoSteps)
	if step.ShouldAutomaticallyAbortOnError() {
		return autoAbort(step, runErr, args)
	}
	args.RunState.RunSteps.Prepend(step.CreateContinueSteps()...)
	err = args.RunState.MarkAsUnfinished(&args.Run.Backend)
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
	message := runErr.Error()
	if !args.RunState.IsAbort && !args.RunState.IsUndo {
		message += messages.AbortContinueGuidance
	}
	if args.RunState.UnfinishedDetails.CanSkip {
		message += messages.ContinueSkipGuidance
	}
	message += "\n"
	return fmt.Errorf(message)
}
