package validate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/persistence"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/runvm"
	"github.com/git-town/git-town/v9/src/undo"
)

// HandleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
// TODO: convert arguments to struct.
func HandleUnfinishedState(run *git.ProdRunner, connector hosting.Connector, rootDir domain.RepoRootDir, lineage config.Lineage, initialBranchesSnapshot domain.BranchesSnapshot, initialConfigSnapshot undo.ConfigSnapshot, initialStashSnapshot domain.StashSnapshot, pushHook bool) (quit bool, err error) {
	runState, err := persistence.Load(rootDir)
	if err != nil {
		return false, fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return false, nil
	}
	response, err := dialog.AskHowToHandleUnfinishedRunState(
		runState.Command,
		runState.UnfinishedDetails.EndBranch,
		runState.UnfinishedDetails.EndTime,
		runState.UnfinishedDetails.CanSkip,
	)
	if err != nil {
		return quit, err
	}
	switch response {
	case dialog.ResponseDiscard:
		return discardRunstate(rootDir)
	case dialog.ResponseContinue:
		return continueRunstate(run, runState, connector, rootDir, lineage, initialBranchesSnapshot, initialConfigSnapshot, initialStashSnapshot, pushHook)
	case dialog.ResponseAbort:
		return abortRunstate(run, runState, connector, rootDir, lineage, initialBranchesSnapshot, initialConfigSnapshot, initialStashSnapshot, pushHook)
	case dialog.ResponseSkip:
		return skipRunstate(run, runState, connector, rootDir, lineage, initialBranchesSnapshot, initialConfigSnapshot, initialStashSnapshot, pushHook)
	case dialog.ResponseQuit:
		return true, nil
	default:
		return false, fmt.Errorf(messages.DialogUnexpectedResponse, response)
	}
}

func abortRunstate(run *git.ProdRunner, runState *runstate.RunState, connector hosting.Connector, rootDir domain.RepoRootDir, lineage config.Lineage, initialBranchesSnapshot domain.BranchesSnapshot, initialConfigSnapshot undo.ConfigSnapshot, initialStashSnapshot domain.StashSnapshot, pushHook bool) (bool, error) {
	abortRunState := runState.CreateAbortRunState()
	return true, runvm.Execute(runvm.ExecuteArgs{
		RunState:                &abortRunState,
		Run:                     run,
		Connector:               connector,
		RootDir:                 rootDir,
		Lineage:                 lineage,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   initialConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
		NoPushHook:              !pushHook,
	})
}

func continueRunstate(run *git.ProdRunner, runState *runstate.RunState, connector hosting.Connector, rootDir domain.RepoRootDir, lineage config.Lineage, initialBranchesSnapshot domain.BranchesSnapshot, initialConfigSnapshot undo.ConfigSnapshot, initialStashSnapshot domain.StashSnapshot, pushHook bool) (bool, error) {
	hasConflicts, err := run.Backend.HasConflicts()
	if err != nil {
		return false, err
	}
	if hasConflicts {
		return false, fmt.Errorf(messages.ContinueUnresolvedConflicts)
	}
	return true, runvm.Execute(runvm.ExecuteArgs{
		RunState:                runState,
		Run:                     run,
		Connector:               connector,
		Lineage:                 lineage,
		RootDir:                 rootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   initialConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
		NoPushHook:              !pushHook,
	})
}

func discardRunstate(rootDir domain.RepoRootDir) (bool, error) {
	err := persistence.Delete(rootDir)
	return false, err
}

func skipRunstate(run *git.ProdRunner, runState *runstate.RunState, connector hosting.Connector, rootDir domain.RepoRootDir, lineage config.Lineage, initialBranchesSnapshot domain.BranchesSnapshot, initialConfigSnapshot undo.ConfigSnapshot, initialStashSnapshot domain.StashSnapshot, pushHook bool) (bool, error) {
	skipRunState := runState.CreateSkipRunState()
	return true, runvm.Execute(runvm.ExecuteArgs{
		RunState:                &skipRunState,
		Run:                     run,
		Connector:               connector,
		Lineage:                 lineage,
		RootDir:                 rootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   initialConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
		NoPushHook:              !pushHook,
	})
}
