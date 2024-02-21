package interpreter

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/print"
	"github.com/git-town/git-town/v12/src/config/gitconfig"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	"github.com/git-town/git-town/v12/src/vm/shared"
	"github.com/git-town/git-town/v12/src/vm/statefile"
)

// errored is called when the given opcode has resulted in the given error.
func errored(failedOpcode shared.Opcode, runErr error, args ExecuteArgs) error {
	var err error
	args.RunState.EndBranchesSnapshot, err = args.Run.Backend.BranchesSnapshot()
	if err != nil {
		return err
	}
	configGitAccess := gitconfig.Access{Runner: args.Run.Backend.Runner}
	globalSnapshot, _, err := configGitAccess.LoadGlobal()
	if err != nil {
		return err
	}
	localSnapshot, _, err := configGitAccess.LoadLocal()
	if err != nil {
		return err
	}
	args.RunState.EndConfigSnapshot = undoconfig.ConfigSnapshot{
		Global: globalSnapshot,
		Local:  localSnapshot,
	}
	args.RunState.EndStashSize, err = args.Run.Backend.StashSize()
	if err != nil {
		return err
	}
	args.RunState.AbortProgram.Add(failedOpcode.CreateAbortProgram()...)
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
	if args.RunState.Command == "sync" && !(repoStatus.RebaseInProgress && args.Run.Config.FullConfig.IsMainBranch(currentBranch)) {
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
	return errors.New(message)
}
