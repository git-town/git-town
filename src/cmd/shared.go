package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
)

func validateIsConfigured(repo *git.ProdRepo) error {
	err := dialog.EnsureIsConfigured(repo)
	if err != nil {
		return err
	}
	return repo.RemoveOutdatedConfiguration()
}

// ValidateIsRepository asserts that the current directory is in a Git repository.
// If so, it also navigates to the root directory.
func ValidateIsRepository(repo *git.ProdRepo) error {
	if repo.Silent.IsRepository() {
		return repo.NavigateToRootIfNecessary()
	}
	return errors.New("this is not a Git repository")
}

func appendStepList(config *appendConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	for _, branch := range append(config.ancestorBranches, config.parentBranch) {
		steps, err := syncBranchSteps(branch, true, repo)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	}
	result.Append(&steps.CreateBranchStep{Branch: config.targetBranch, StartingPoint: config.parentBranch})
	result.Append(&steps.SetParentBranchStep{Branch: config.targetBranch, ParentBranch: config.parentBranch})
	result.Append(&steps.CheckoutBranchStep{Branch: config.targetBranch})
	if config.hasOrigin && config.shouldNewBranchPush && !config.isOffline {
		result.Append(&steps.CreateTrackingBranchStep{Branch: config.targetBranch, NoPushHook: config.noPushHook})
	}
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	return result, err
}

// handleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
//
//nolint:nonamedreturns  // return value isn't obvious from function name
func handleUnfinishedState(repo *git.ProdRepo, connector hosting.Connector) (quit bool, err error) {
	runState, err := runstate.Load(repo)
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
		err = runstate.Delete(repo)
		return false, err
	case dialog.ResponseTypeContinue:
		hasConflicts, err := repo.Silent.HasConflicts()
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
	default:
		return false, fmt.Errorf("unknown response: %s", response)
	}
}
