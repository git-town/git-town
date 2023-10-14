package runvm

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/step"
	"github.com/git-town/git-town/v9/src/undo"
	"github.com/git-town/git-town/v9/src/vm/persistence"
)

// errored is called when the given step has resulted in the given error.
func errored(failedStep step.Step, runErr error, args ExecuteArgs) error {
	args.RunState.AbortSteps.Add(failedStep.CreateAbortSteps()...)
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
	args.RunState.UndoSteps.AddList(undoSteps)
	if failedStep.ShouldAutomaticallyAbortOnError() {
		return autoAbort(failedStep, runErr, args)
	}
	args.RunState.RunSteps.Prepend(failedStep.CreateContinueSteps()...)
	err = args.RunState.MarkAsUnfinished(&args.Run.Backend)
	if err != nil {
		return err
	}
	currentBranch, err := args.Run.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	repoStatus, err := args.Run.Backend.RepoStatus()
	if err != nil {
		return err
	}
	if args.RunState.Command == "sync" && !(repoStatus.RebaseInProgress && args.Run.Config.IsMainBranch(currentBranch)) {
		args.RunState.UnfinishedDetails.CanSkip = true
	}
	err = persistence.Save(args.RunState, args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	PrintFooter(args.Debug, args.Run.CommandsCounter.Count(), args.Run.FinalMessages.Result())
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
