package cmd

import (
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
)

// These variables are set at build time.
var (
	version,
	buildDate string
)

const dryRunFlagDescription = "Print the commands but don't run them"

func validateIsConfigured() error {
	prompt.EnsureIsConfigured()
	git.Config().RemoveOutdatedConfiguration()
	return nil
}

func ensureIsNotInUnfinishedState(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	runState, err := steps.LoadPreviousRunState()
	if err != nil {
		fmt.Printf("cannot load previous run state: %v\n", err)
		os.Exit(1)
	}
	if runState != nil && runState.IsUnfinished() {
		response := prompt.AskHowToHandleUnfinishedRunState(
			runState.Command,
			runState.UnfinishedDetails.EndBranch,
			runState.UnfinishedDetails.EndTime,
			runState.UnfinishedDetails.CanSkip,
		)
		switch response {
		case prompt.ResponseTypeDiscard:
			return steps.DeletePreviousRunState()
		case prompt.ResponseTypeContinue:
			if git.HasConflicts() {
				return fmt.Errorf("you must resolve the conflicts before continuing")
			}
			err = steps.Run(runState, repo, driver)
		case prompt.ResponseTypeAbort:
			abortRunState := runState.CreateAbortRunState()
			err = steps.Run(&abortRunState, repo, driver)
		case prompt.ResponseTypeSkip:
			skipRunState := runState.CreateSkipRunState()
			err = steps.Run(&skipRunState, repo, driver)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
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
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	return result, nil
}
