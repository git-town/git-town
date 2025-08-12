package fullinterpreter

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/state/runlog"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// errored is called when the given opcode has resulted in the given error.
func errored(failedOpcode shared.Opcode, runErr error, args ExecuteArgs) error {
	endBranchesSnapshot, err := args.Git.BranchesSnapshot(args.Backend)
	if err != nil {
		return err
	}
	err = runlog.Write(runlog.EventEnd, endBranchesSnapshot.Branches, Some(args.RunState.Command), args.RootDir)
	if err != nil {
		return err
	}
	args.RunState.EndBranchesSnapshot = Some(endBranchesSnapshot)
	globalSnapshot, err := gitconfig.LoadSnapshot(args.Backend, Some(configdomain.ConfigScopeGlobal), configdomain.UpdateOutdatedNo)
	if err != nil {
		return err
	}
	localSnapshot, err := gitconfig.LoadSnapshot(args.Backend, Some(configdomain.ConfigScopeLocal), configdomain.UpdateOutdatedNo)
	if err != nil {
		return err
	}
	args.RunState.EndConfigSnapshot = Some(undoconfig.EndConfigSnapshot{
		Global: globalSnapshot,
		Local:  localSnapshot,
	})
	endStashSize, err := args.Git.StashSize(args.Backend)
	if err != nil {
		return err
	}
	args.RunState.EndStashSize = Some(endStashSize)
	args.RunState.AbortProgram.Add(failedOpcode.Abort()...)
	if failedOpcode.ShouldUndoOnError() {
		return autoUndo(failedOpcode, runErr, args)
	}
	continueProgram := failedOpcode.Continue()
	if len(continueProgram) == 0 {
		continueProgram = []shared.Opcode{failedOpcode}
	}
	args.RunState.RunProgram.Prepend(continueProgram...)
	currentBranch, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	repoStatus, err := args.Git.RepoStatus(args.Backend)
	if err != nil {
		return err
	}
	canSkip := false
	if args.RunState.Command == "sync" && !(repoStatus.RebaseInProgress && args.Config.ValidatedConfigData.IsMainBranch(currentBranch)) {
		canSkip = true
	}
	if args.RunState.Command == "walk" {
		canSkip = true
	}
	if err = args.RunState.MarkAsUnfinished(args.Git, args.Backend, canSkip); err != nil {
		return err
	}
	if err = runstate.Save(args.RunState, args.RootDir); err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	print.Footer(args.Config.NormalConfig.Verbose, args.CommandsCounter.Immutable(), args.FinalMessages.Result())
	message := runErr.Error()
	message += messages.UndoContinueGuidance
	if unfinishedDetails, hasUnfinishedDetails := args.RunState.UnfinishedDetails.Get(); hasUnfinishedDetails {
		if unfinishedDetails.CanSkip {
			message += messages.ContinueSkipGuidance
		}
	}
	message += "\n"
	args.Inputs.VerifyAllUsed()
	return errors.New(message)
}
