package validate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
)

// HandleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
func HandleUnfinishedState(run *git.ProdRunner, connector hosting.Connector, rootDir string) (quit bool, err error) {
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
		err = runstate.Delete(rootDir)
		return false, err
	case dialog.ResponseContinue:
		hasConflicts, err := run.Backend.HasConflicts()
		if err != nil {
			return false, err
		}
		if hasConflicts {
			return false, fmt.Errorf(messages.ContinueUnresolvedConflicts)
		}
		return true, runstate.Execute(runstate.ExecuteArgs{
			RunState:  runState,
			Run:       run,
			Connector: connector,
			RootDir:   rootDir,
		})
	case dialog.ResponseAbort:
		abortRunState := runState.CreateAbortRunState()
		return true, runstate.Execute(runstate.ExecuteArgs{
			RunState:  &abortRunState,
			Run:       run,
			Connector: connector,
			RootDir:   rootDir,
		})
	case dialog.ResponseSkip:
		skipRunState := runState.CreateSkipRunState()
		return true, runstate.Execute(runstate.ExecuteArgs{
			RunState:  &skipRunState,
			Run:       run,
			Connector: connector,
			RootDir:   rootDir,
		})
	case dialog.ResponseQuit:
		return true, nil
	default:
		return false, fmt.Errorf(messages.DialogUnexpectedResponse, response)
	}
}
