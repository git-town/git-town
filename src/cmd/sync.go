package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/git-town/git-town/v7/src/stringslice"
	"github.com/spf13/cobra"
)

type syncConfig struct {
	branchesToSync                           []string
	hasOrigin                                bool
	initialBranch                            string
	isOffline                                bool
	localBranchesWithDeletedTrackingBranches []string
	mainBranch                               string
	shouldPushTags                           bool
}

func syncCmd(repo *git.ProdRepo) *cobra.Command {
	var allFlag bool
	var dryRunFlag bool
	syncCmd := cobra.Command{
		Use:   "sync",
		Short: "Updates the current branch with all relevant changes",
		Long: fmt.Sprintf(`Updates the current branch with all relevant changes

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
You can disable this by running "git config %s false".`, config.SyncUpstream),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := createSyncConfig(allFlag, repo)
			if err != nil {
				cli.Exit(err)
			}
			stepList, err := createSyncStepList(config, repo)
			if err != nil {
				cli.Exit(err)
			}
			runState := runstate.New("sync", stepList)
			err = runstate.Execute(runState, repo, nil)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			if dryRunFlag {
				currentBranch, err := repo.Silent.CurrentBranch()
				if err != nil {
					return err
				}
				repo.DryRun.Activate(currentBranch)
			}
			if err := validateIsConfigured(repo); err != nil {
				return err
			}
			exit, err := handleUnfinishedState(repo, nil)
			if err != nil {
				return err
			}
			if exit {
				os.Exit(0)
			}
			return nil
		},
	}
	syncCmd.Flags().BoolVar(&allFlag, "all", false, "Sync all local branches")
	syncCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Print the commands but don't run them")
	return &syncCmd
}

func createSyncConfig(allFlag bool, repo *git.ProdRepo) (syncConfig, error) {
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return syncConfig{}, err
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return syncConfig{}, err
	}
	result := syncConfig{
		hasOrigin:  hasOrigin,
		isOffline:  isOffline,
		mainBranch: repo.Config.MainBranch(),
	}
	if result.hasOrigin && !result.isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return syncConfig{}, err
		}
	}
	result.localBranchesWithDeletedTrackingBranches, err = repo.Silent.LocalBranchesWithDeletedTrackingBranches()
	if err != nil {
		return syncConfig{}, err
	}
	result.initialBranch, err = repo.Silent.CurrentBranch()
	if err != nil {
		return syncConfig{}, err
	}
	parentDialog := dialog.ParentBranches{}
	if allFlag {
		branches, err := repo.Silent.LocalBranchesMainFirst()
		if err != nil {
			return syncConfig{}, err
		}
		err = parentDialog.EnsureKnowsParentBranches(branches, repo)
		if err != nil {
			return syncConfig{}, err
		}
		result.branchesToSync = branches
		result.shouldPushTags = true
	} else {
		err = parentDialog.EnsureKnowsParentBranches([]string{result.initialBranch}, repo)
		if err != nil {
			return syncConfig{}, err
		}
		result.branchesToSync = append(repo.Config.AncestorBranches(result.initialBranch), result.initialBranch)
		result.shouldPushTags = !repo.Config.IsFeatureBranch(result.initialBranch)
	}
	return result, nil
}

func createSyncStepList(config syncConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	for _, branchName := range config.branchesToSync {
		var stepsForBranch runstate.StepList
		var err error
		if stringslice.Contains(config.localBranchesWithDeletedTrackingBranches, branchName) {
			stepsForBranch, err = createDeleteBranchStepList(branchName, config, repo)
		} else {
			stepsForBranch, err = runstate.SyncBranchSteps(branchName, true, repo)
		}
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(stepsForBranch)
	}
	var finalBranch string
	if stringslice.Contains(config.localBranchesWithDeletedTrackingBranches, config.initialBranch) {
		finalBranch = config.mainBranch
	} else {
		finalBranch = config.initialBranch
	}
	result.Append(&steps.CheckoutBranchStep{BranchName: finalBranch})
	if config.hasOrigin && config.shouldPushTags && !config.isOffline {
		result.Append(&steps.PushTagsStep{})
	}
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	return result, err
}

func createDeleteBranchStepList(branchName string, config syncConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	if config.initialBranch == branchName {
		result.Append(&steps.CheckoutBranchStep{BranchName: config.mainBranch})
	}
	parent := repo.Config.ParentBranch(branchName)
	if parent != "" {
		for _, child := range repo.Config.ChildBranches(branchName) {
			result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: parent})
		}
		result.Append(&steps.DeleteParentBranchStep{BranchName: branchName})
	}
	if repo.Config.IsPerennialBranch(branchName) {
		result.Append(&steps.RemoveFromPerennialBranchesStep{BranchName: branchName})
	}
	result.Append(&steps.DeleteLocalBranchStep{BranchName: branchName})
	return result, nil
}
