package fullinterpreter

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v20/internal/cli/print"
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/config/gitconfig"
	"github.com/git-town/git-town/v20/internal/messages"
	"github.com/git-town/git-town/v20/internal/undo/undoconfig"
	"github.com/git-town/git-town/v20/internal/vm/shared"
	"github.com/git-town/git-town/v20/internal/vm/statefile"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// errored is called when the given opcode has resulted in the given error.
func errored(failedOpcode shared.Opcode, runErr error, args ExecuteArgs) error {
	endBranchesSnapshot, err := args.Git.BranchesSnapshot(args.Backend)
	if err != nil {
		return err
	}
	args.RunState.EndBranchesSnapshot = Some(endBranchesSnapshot)
	configGitAccess := gitconfig.Access{Runner: args.Backend}
	globalSnapshot, err := configGitAccess.Load(Some(configdomain.ConfigScopeGlobal), configdomain.UpdateOutdatedNo)
	if err != nil {
		return err
	}
	localSnapshot, err := configGitAccess.Load(Some(configdomain.ConfigScopeLocal), configdomain.UpdateOutdatedNo)
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
	args.RunState.AbortProgram.Add(failedOpcode.AbortProgram()...)
	if failedOpcode.ShouldUndoOnError() {
		return autoUndo(failedOpcode, runErr, args)
	}
	continueProgram := failedOpcode.ContinueProgram()
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
	err = args.RunState.MarkAsUnfinished(args.Git, args.Backend, canSkip)
	if err != nil {
		return err
	}
	err = statefile.Save(args.RunState, args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	print.Footer(args.Verbose, args.CommandsCounter.Immutable(), args.FinalMessages.Result())
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
