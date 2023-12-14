package interpreter

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/print"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/git-town/git-town/v11/src/undo"
	"github.com/git-town/git-town/v11/src/vm/shared"
	"github.com/git-town/git-town/v11/src/vm/statefile"
)

// errored is called when the given opcode has resulted in the given error.
func errored(failedOpcode shared.Opcode, runErr error, args ExecuteArgs) error {
	args.RunState.AbortProgram.Add(failedOpcode.CreateAbortProgram()...)
	undoProgram, err := undo.CreateUndoProgram(undo.CreateUndoProgramArgs{
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
	args.RunState.UndoProgram.AddProgram(undoProgram)
	if failedOpcode.ShouldAutomaticallyUndoOnError() {
		return autoUndo(failedOpcode, runErr, args)
	}
	args.RunState.RunProgram.Prepend(failedOpcode.CreateContinueProgram()...)
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
	if args.RunState.Command == "sync" && !(repoStatus.RebaseInProgress && args.Run.GitTown.IsMainBranch(currentBranch)) {
		args.RunState.UnfinishedDetails.CanSkip = true
	}
	err = statefile.Save(args.RunState, args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	print.Footer(args.Verbose, args.Run.CommandsCounter.Count(), args.Run.FinalMessages.Result())
	message := runErr.Error()
	if !args.RunState.IsUndo {
		message += messages.UndoContinueGuidance
	}
	if args.RunState.UnfinishedDetails.CanSkip {
		message += messages.ContinueSkipGuidance
	}
	message += "\n"
	return fmt.Errorf(message)
}
