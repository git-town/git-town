package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/git-town/git-town/v7/src/userinput"
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
	err := userinput.EnsureIsConfigured(repo)
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

func createAppendStepList(config appendConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	for _, branchName := range append(config.ancestorBranches, config.parentBranch) {
		steps, err := runstate.SyncBranchSteps(branchName, true, repo)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	}
	result.Append(&steps.CreateBranchStep{BranchName: config.targetBranch, StartingPoint: config.parentBranch})
	result.Append(&steps.SetParentBranchStep{BranchName: config.targetBranch, ParentBranchName: config.parentBranch})
	result.Append(&steps.CheckoutBranchStep{BranchName: config.targetBranch})
	if config.hasOrigin && config.shouldNewBranchPush && !config.isOffline {
		result.Append(&steps.CreateTrackingBranchStep{BranchName: config.targetBranch})
	}
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	return result, err
}

// handleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
//
//nolint:nonamedreturns return value isn't obvious from function name
func handleUnfinishedState(repo *git.ProdRepo, driver hosting.Driver) (quit bool, err error) {
	runState, err := runstate.Load(repo)
	if err != nil {
		return false, fmt.Errorf("cannot load previous run state: %w", err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return false, nil
	}
	response, err := userinput.AskHowToHandleUnfinishedRunState(
		runState.Command,
		runState.UnfinishedDetails.EndBranch,
		runState.UnfinishedDetails.EndTime,
		runState.UnfinishedDetails.CanSkip,
	)
	if err != nil {
		return quit, err
	}
	switch response {
	case userinput.ResponseTypeDiscard:
		err = runstate.Delete(repo)
		return false, err
	case userinput.ResponseTypeContinue:
		hasConflicts, err := repo.Silent.HasConflicts()
		if err != nil {
			return false, err
		}
		if hasConflicts {
			return false, fmt.Errorf("you must resolve the conflicts before continuing")
		}
		return true, runstate.Execute(runState, repo, driver)
	case userinput.ResponseTypeAbort:
		abortRunState := runState.CreateAbortRunState()
		return true, runstate.Execute(&abortRunState, repo, driver)
	case userinput.ResponseTypeSkip:
		skipRunState := runState.CreateSkipRunState()
		return true, runstate.Execute(&skipRunState, repo, driver)
	default:
		return false, fmt.Errorf("unknown response: %s", response)
	}
}
