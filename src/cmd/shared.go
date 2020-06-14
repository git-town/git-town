package cmd

import (
	"errors"
	"fmt"
	"os"

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
	version,
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
func ValidateIsRepository(repo *git.ProdRepo) error {
	if repo.Silent.IsRepository() {
		return nil
	}
	return errors.New("this is not a Git repository")
}

func ensureIsNotInUnfinishedState(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	runState, err := steps.LoadPreviousRunState(repo)
	if err != nil {
		return fmt.Errorf("cannot load previous run state: %w", err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return nil
	}
	response := prompt.AskHowToHandleUnfinishedRunState(
		runState.Command,
		runState.UnfinishedDetails.EndBranch,
		runState.UnfinishedDetails.EndTime,
		runState.UnfinishedDetails.CanSkip,
	)
	switch response {
	case prompt.ResponseTypeDiscard:
		return steps.DeletePreviousRunState(repo)
	case prompt.ResponseTypeContinue:
		hasConflicts, err := repo.Silent.HasConflicts()
		if err != nil {
			return err
		}
		if hasConflicts {
			return fmt.Errorf("you must resolve the conflicts before continuing")
		}
		err = steps.Run(runState, repo, driver)
		if err != nil {
			return err
		}
	case prompt.ResponseTypeAbort:
		abortRunState := runState.CreateAbortRunState()
		err = steps.Run(&abortRunState, repo, driver)
	case prompt.ResponseTypeSkip:
		skipRunState := runState.CreateSkipRunState()
		err = steps.Run(&skipRunState, repo, driver)
	default:
		return fmt.Errorf("unknown response: %s", response)
	}
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
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
