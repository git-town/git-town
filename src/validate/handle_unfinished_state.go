package validate

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
)

// HandleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
//
//nolint:nonamedreturns  // return value isn't obvious from function name
func HandleUnfinishedState(repo *git.ProdRepo, connector hosting.Connector) (quit bool, err error) {
	runState, err := runstate.Load(&repo.Backend)
	if err != nil {
		return false, fmt.Errorf("cannot load previous run state: %w", err)
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
	case dialog.ResponseTypeDiscard:
		err = runstate.Delete(&repo.Backend)
		return false, err
	case dialog.ResponseTypeContinue:
		hasConflicts, err := repo.Backend.HasConflicts()
		if err != nil {
			return false, err
		}
		if hasConflicts {
			return false, fmt.Errorf("you must resolve the conflicts before continuing")
		}
		return true, runstate.Execute(runState, repo, connector)
	case dialog.ResponseTypeAbort:
		abortRunState := runState.CreateAbortRunState()
		return true, runstate.Execute(&abortRunState, repo, connector)
	case dialog.ResponseTypeSkip:
		skipRunState := runState.CreateSkipRunState()
		return true, runstate.Execute(&skipRunState, repo, connector)
	case dialog.ResponseTypeQuit:
		return true, nil
	default:
		return false, fmt.Errorf("unknown response: %s", response)
	}
}
