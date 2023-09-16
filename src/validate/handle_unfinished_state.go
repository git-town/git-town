package validate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
)

// HandleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
func HandleUnfinishedState(run *git.ProdRunner, connector hosting.Connector, rootDir string, lineage config.Lineage, initialBranchesSnapshot runstate.BranchesSnapshot, initialConfigSnapshot runstate.ConfigSnapshot) (quit bool, err error) {
	runState, err := runstate.Load(rootDir)
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
		return continueRunstate(run, runState, connector, rootDir, lineage, initialBranchesSnapshot, initialConfigSnapshot)
	case dialog.ResponseAbort:
		return abortRunstate(run, runState, connector, rootDir, lineage, initialBranchesSnapshot, initialConfigSnapshot)
	case dialog.ResponseSkip:
		return skipRunstate(run, runState, connector, rootDir, lineage, initialBranchesSnapshot, initialConfigSnapshot)
	case dialog.ResponseQuit:
		return true, nil
	default:
		return false, fmt.Errorf(messages.DialogUnexpectedResponse, response)
	}
}

func abortRunstate(run *git.ProdRunner, runState *runstate.RunState, connector hosting.Connector, rootDir string, lineage config.Lineage, initialBranchesSnapshot runstate.BranchesSnapshot, initialConfigSnapshot runstate.ConfigSnapshot) (bool, error) {
	abortRunState := runState.CreateAbortRunState()
	return true, runstate.Execute(runstate.ExecuteArgs{
		RunState:                &abortRunState,
		Run:                     run,
		Connector:               connector,
		RootDir:                 rootDir,
		Lineage:                 lineage,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   initialConfigSnapshot,
	})
}

func continueRunstate(run *git.ProdRunner, runState *runstate.RunState, connector hosting.Connector, rootDir string, lineage config.Lineage, initialBranchesSnapshot runstate.BranchesSnapshot, initialConfigSnapshot runstate.ConfigSnapshot) (bool, error) {
	hasConflicts, err := run.Backend.HasConflicts()
	if err != nil {
		return false, err
	}
	if hasConflicts {
		return false, fmt.Errorf(messages.ContinueUnresolvedConflicts)
	}
	return true, runstate.Execute(runstate.ExecuteArgs{
		RunState:                runState,
		Run:                     run,
		Connector:               connector,
		Lineage:                 lineage,
		RootDir:                 rootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   initialConfigSnapshot,
	})
}

func discardRunstate(rootDir string) (bool, error) {
	err := runstate.Delete(rootDir)
	return false, err
}

func skipRunstate(run *git.ProdRunner, runState *runstate.RunState, connector hosting.Connector, rootDir string, lineage config.Lineage, initialBranchesSnapshot runstate.BranchesSnapshot, initialConfigSnapshot runstate.ConfigSnapshot) (bool, error) {
	skipRunState := runState.CreateSkipRunState()
	return true, runstate.Execute(runstate.ExecuteArgs{
		RunState:                &skipRunState,
		Run:                     run,
		Connector:               connector,
		Lineage:                 lineage,
		RootDir:                 rootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   initialConfigSnapshot,
	})
}
