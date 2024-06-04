package interpreter

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/vm/shared"
	"github.com/git-town/git-town/v14/src/vm/statefile"
)

// errored is called when the given opcode has resulted in the given error.
func errored(failedOpcode shared.Opcode, runErr error, args ExecuteArgs) error {
	endBranchesSnapshot, err := args.Git.BranchesSnapshot(args.Backend)
	if err != nil {
		return err
	}
	args.RunState.EndBranchesSnapshot = Some(endBranchesSnapshot)
	configGitAccess := gitconfig.Access{Runner: args.Backend}
	globalSnapshot, _, err := configGitAccess.LoadGlobal(false)
	if err != nil {
		return err
	}
	localSnapshot, _, err := configGitAccess.LoadLocal(false)
	if err != nil {
		return err
	}
	args.RunState.EndConfigSnapshot = Some(undoconfig.ConfigSnapshot{
		Global: globalSnapshot,
		Local:  localSnapshot,
	})
	endStashSize, err := args.Git.StashSize(args.Backend)
	if err != nil {
		return err
	}
	args.RunState.EndStashSize = Some(endStashSize)
	args.RunState.AbortProgram.Add(failedOpcode.CreateAbortProgram()...)
	if failedOpcode.ShouldAutomaticallyUndoOnError() {
		return autoUndo(failedOpcode, runErr, args)
	}
	args.RunState.RunProgram.Prepend(failedOpcode.CreateContinueProgram()...)
	err = args.RunState.MarkAsUnfinished(args.Git, args.Backend)
	if err != nil {
		return err
	}
	currentBranch, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	repoStatus, err := args.Git.RepoStatus(args.Backend)
	if err != nil {
		return err
	}
	if args.RunState.Command == "sync" && !(repoStatus.RebaseInProgress && args.Config.Config.IsMainBranch(currentBranch)) {
		if unfinishedDetails, hasUnfinishedDetails := args.RunState.UnfinishedDetails.Get(); hasUnfinishedDetails {
			unfinishedDetails.CanSkip = true
		}
	}
	err = statefile.Save(args.RunState, args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	print.Footer(args.Verbose, args.CommandsCounter.Count(), args.FinalMessages.Result())
	message := runErr.Error()
	message += messages.UndoContinueGuidance
	if unfinishedDetails, hasUnfinishedDetails := args.RunState.UnfinishedDetails.Get(); hasUnfinishedDetails {
		if unfinishedDetails.CanSkip {
			message += messages.ContinueSkipGuidance
		}
	}
	message += "\n"
	return errors.New(message)
}
