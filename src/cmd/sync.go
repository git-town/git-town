package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/dryrun"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/steps"

	"github.com/spf13/cobra"
)

type syncConfig struct {
	initialBranch  string
	branchesToSync []string
	shouldPushTags bool
	hasOrigin      bool
	isOffline      bool
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Updates the current branch with all relevant changes",
	Long: `Updates the current branch with all relevant changes

Synchronizes the current branch with the rest of the world.

When run on a feature branch
- syncs all ancestor branches
- pulls updates for the current branch
- merges the parent branch into the current branch
- pushes the current branch

When run on the main branch or a perennial branch
- pulls and pushes updates for the current branch
- pushes tags

If the repository contains an "upstream" remote,
syncs the main branch with its upstream counterpart.
You can disable this by running "git config git-town.sync-upstream false".`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := getSyncConfig(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		stepList, err := getSyncStepList(config, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		runState := steps.NewRunState("sync", stepList)
		err = steps.Run(runState, prodRepo, nil)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := ValidateIsRepository(prodRepo); err != nil {
			return err
		}
		if dryRunFlag {
			currentBranch, err := prodRepo.Silent.CurrentBranch()
			if err != nil {
				return err
			}
			dryrun.Activate(currentBranch)
		}
		if err := validateIsConfigured(prodRepo); err != nil {
			return err
		}
		_, err := isInUnfinishedState(prodRepo, nil)
		return err
	},
}

func getSyncConfig(repo *git.ProdRepo) (result syncConfig, err error) {
	result.hasOrigin, err = repo.Silent.HasRemote("origin")
	if err != nil {
		return result, err
	}
	result.isOffline = git.Config().IsOffline()
	if result.hasOrigin && !result.isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return result, err
		}
	}
	result.initialBranch, err = repo.Silent.CurrentBranch()
	if err != nil {
		return result, err
	}
	if allFlag {
		branches, err := repo.Silent.LocalBranchesMainFirst()
		if err != nil {
			return result, err
		}
		err = prompt.EnsureKnowsParentBranches(branches, repo)
		if err != nil {
			return result, err
		}
		result.branchesToSync = branches
		result.shouldPushTags = true
	} else {
		err = prompt.EnsureKnowsParentBranches([]string{result.initialBranch}, repo)
		if err != nil {
			return result, err
		}
		result.branchesToSync = append(git.Config().GetAncestorBranches(result.initialBranch), result.initialBranch)
		result.shouldPushTags = !git.Config().IsFeatureBranch(result.initialBranch)
	}
	return result, nil
}

func getSyncStepList(config syncConfig, repo *git.ProdRepo) (result steps.StepList, err error) {
	for _, branchName := range config.branchesToSync {
		steps, err := steps.GetSyncBranchSteps(branchName, true, repo)
		if err != nil {
			return result, err
		}
		result.AppendList(steps)
	}
	result.Append(&steps.CheckoutBranchStep{BranchName: config.initialBranch})
	if config.hasOrigin && config.shouldPushTags && !config.isOffline {
		result.Append(&steps.PushTagsStep{})
	}
	err = result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	return result, err
}

func isInUnfinishedState(repo *git.ProdRepo, driver drivers.CodeHostingDriver) (quit bool, err error) {
	runState, err := steps.LoadPreviousRunState(repo)
	if err != nil {
		return false, fmt.Errorf("cannot load previous run state: %w", err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return false, nil
	}
	response := prompt.AskHowToHandleUnfinishedRunState(
		runState.Command,
		runState.UnfinishedDetails.EndBranch,
		runState.UnfinishedDetails.EndTime,
		runState.UnfinishedDetails.CanSkip,
	)
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
		err = steps.Run(runState, repo, driver)
		if err != nil {
			return false, err
		}
		os.Exit(0)
	case prompt.ResponseTypeAbort:
		abortRunState := runState.CreateAbortRunState()
		err = steps.Run(&abortRunState, repo, driver)
		if err != nil {
			return false, err
		}
		os.Exit(0)
	case prompt.ResponseTypeSkip:
		skipRunState := runState.CreateSkipRunState()
		err = steps.Run(&skipRunState, repo, driver)
		if err != nil {
			return false, err
		}
		os.Exit(0)
	default:
		return true, fmt.Errorf("unknown response: %s", response)
	}
	return false, nil
}

func init() {
	syncCmd.Flags().BoolVar(&allFlag, "all", false, "Sync all local branches")
	syncCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, dryRunFlagDescription)
	RootCmd.AddCommand(syncCmd)
}
