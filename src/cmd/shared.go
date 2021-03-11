package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/steps"
)

// These variables represent command-line flags.
var (
	allFlag,
	debugFlag,
	dryRunFlag,
	globalFlag bool
	prodRepo = git.NewProdRepo()
)

// These variables are set at build time.
var (
	version   string
	buildDate string
)

const dryRunFlagDescription = "Print the commands but don't run them"

func validateIsConfigured(repo *git.ProdRepo) error {
	err := prompt.EnsureIsConfigured(repo)
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

func getAppendStepList(config appendConfig, repo *git.ProdRepo) (result steps.StepList, err error) {
	for _, branchName := range append(config.ancestorBranches, config.parentBranch) {
		steps, err := steps.GetSyncBranchSteps(branchName, true, repo)
		if err != nil {
			return result, err
		}
		result.AppendList(steps)
	}
	result.Append(&steps.CreateBranchStep{BranchName: config.targetBranch, StartingPoint: config.parentBranch})
	result.Append(&steps.SetParentBranchStep{BranchName: config.targetBranch, ParentBranchName: config.parentBranch})
	result.Append(&steps.CheckoutBranchStep{BranchName: config.targetBranch})
	if config.hasOrigin && config.shouldNewBranchPush && !config.isOffline {
		result.Append(&steps.CreateTrackingBranchStep{BranchName: config.targetBranch})
	}
	err = result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	return result, err
}

// handleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
func handleUnfinishedState(repo *git.ProdRepo, driver drivers.CodeHostingDriver) (quit bool, err error) {
	runState, err := steps.LoadPreviousRunState(repo)
	if err != nil {
		return false, fmt.Errorf("cannot load previous run state: %w", err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return false, nil
	}
	response, err := prompt.AskHowToHandleUnfinishedRunState(
		runState.Command,
		runState.UnfinishedDetails.EndBranch,
		runState.UnfinishedDetails.EndTime,
		runState.UnfinishedDetails.CanSkip,
	)
	if err != nil {
		return quit, err
	}
	switch response {
	case prompt.ResponseTypeDiscard:
		err = steps.DeletePreviousRunState(repo)
		return false, err
	case prompt.ResponseTypeContinue:
		hasConflicts, err := repo.Silent.HasConflicts()
		if err != nil {
			return false, err
		}
		if hasConflicts {
			return false, fmt.Errorf("you must resolve the conflicts before continuing")
		}
		return true, steps.Run(runState, repo, driver)
	case prompt.ResponseTypeAbort:
		abortRunState := runState.CreateAbortRunState()
		return true, steps.Run(&abortRunState, repo, driver)
	case prompt.ResponseTypeSkip:
		skipRunState := runState.CreateSkipRunState()
		return true, steps.Run(&skipRunState, repo, driver)
	default:
		return false, fmt.Errorf("unknown response: %s", response)
	}
}
