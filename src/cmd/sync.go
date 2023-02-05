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
	branchesToSync            []string
	branchesWithDeletedRemote []string // local branches whose tracking branches have been deleted
	hasOrigin                 bool
	initialBranch             string
	isOffline                 bool
	mainBranch                string
	shouldPushTags            bool
}

func (sc *syncConfig) hasDeletedTrackingBranch(branch string) bool {
	return stringslice.Contains(sc.branchesWithDeletedRemote, branch)
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
			stepList, err := syncBranchesStepList(config, repo)
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
	result.branchesWithDeletedRemote, err = repo.Silent.LocalBranchesWithDeletedTrackingBranches()
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

func syncBranchesStepList(config syncConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	for _, branch := range config.branchesToSync {
		stepsForBranch, err := syncStepsForBranch(branch, config, repo)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(stepsForBranch)
	}
	result.Append(&steps.CheckoutBranchStep{BranchName: finalBranch(config)})
	if config.hasOrigin && config.shouldPushTags && !config.isOffline {
		result.Append(&steps.PushTagsStep{})
	}
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	return result, err
}

func syncStepsForBranch(branch string, config syncConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	if config.hasDeletedTrackingBranch(branch) {
		return deleteBranchSteps(branch, config, repo)
	} else {
		return updateBranchSteps(branch, true, config.branchesWithDeletedRemote, repo)
	}
}

func deleteBranchSteps(branch string, config syncConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	if config.initialBranch == branch {
		result.Append(&steps.CheckoutBranchStep{BranchName: config.mainBranch})
	}
	parent := repo.Config.ParentBranch(branch)
	if parent != "" {
		for _, child := range repo.Config.ChildBranches(branch) {
			result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: parent})
		}
		result.Append(&steps.DeleteParentBranchStep{BranchName: branch})
	}
	if repo.Config.IsPerennialBranch(branch) {
		result.Append(&steps.RemoveFromPerennialBranchesStep{BranchName: branch})
	}
	result.Append(&steps.DeleteLocalBranchStep{BranchName: branch})
	return result, nil
}

//nolint:nestif
func updateBranchSteps(branchName string, pushBranch bool, branchesWithDeletedRemote []string, repo *git.ProdRepo) (runstate.StepList, error) {
	isFeatureBranch := repo.Config.IsFeatureBranch(branchName)
	syncStrategy := repo.Config.SyncStrategy()
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return runstate.StepList{}, err
	}
	pushHook, err := repo.Config.PushHook()
	if err != nil {
		return runstate.StepList{}, err
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return runstate.StepList{}, err
	}
	result := runstate.StepList{}
	if !hasOrigin && !isFeatureBranch {
		return runstate.StepList{}, nil
	}
	result.Append(&steps.CheckoutBranchStep{BranchName: branchName})
	if isFeatureBranch {
		steps, err := syncFeatureBranchSteps(branchName, branchesWithDeletedRemote, repo)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	} else {
		steps, err := syncNonFeatureBranchSteps(branchName, repo)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	}
	if pushBranch && hasOrigin && !isOffline {
		hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branchName)
		if err != nil {
			return runstate.StepList{}, err
		}
		if hasTrackingBranch {
			if isFeatureBranch {
				switch syncStrategy {
				case "merge":
					result.Append(&steps.PushBranchStep{BranchName: branchName, NoPushHook: !pushHook})
				case "rebase":
					result.Append(&steps.PushBranchStep{BranchName: branchName, ForceWithLease: true})
				default:
					return runstate.StepList{}, fmt.Errorf("unknown syncStrategy value: %q", syncStrategy)
				}
			} else {
				result.Append(&steps.PushBranchStep{BranchName: branchName})
			}
		} else {
			result.Append(&steps.CreateTrackingBranchStep{BranchName: branchName})
		}
	}
	return result, nil
}

// Helpers

func syncFeatureBranchSteps(branchName string, branchesWithDeletedRemote []string, repo *git.ProdRepo) (runstate.StepList, error) {
	syncStrategy := repo.Config.SyncStrategy()
	hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branchName)
	if err != nil {
		return runstate.StepList{}, err
	}
	result := runstate.StepList{}
	if hasTrackingBranch {
		switch syncStrategy {
		case "merge":
			result.Append(&steps.MergeBranchStep{BranchName: repo.Silent.TrackingBranchName(branchName)})
		case "rebase":
			result.Append(&steps.RebaseBranchStep{BranchName: repo.Silent.TrackingBranchName(branchName)})
		default:
			return runstate.StepList{}, fmt.Errorf("unknown syncStrategy value: %q", syncStrategy)
		}
	}
	// TODO: the last non-deleted parent branch here
	ancestorBranches := repo.Config.AncestorBranches(branchName)
	ancestorBranches = stringslice.RemoveMany(ancestorBranches, branchesWithDeletedRemote)
	newParentBranch := stringslice.Last(ancestorBranches)
	if newParentBranch == nil {
		return runstate.StepList{}, nil
	}
	switch syncStrategy {
	case "merge":
		result.Append(&steps.MergeBranchStep{BranchName: *newParentBranch})
	case "rebase":
		result.Append(&steps.RebaseBranchStep{BranchName: *newParentBranch})
	default:
		return runstate.StepList{}, fmt.Errorf("unknown syncStrategy value: %q", syncStrategy)
	}
	return result, nil
}

func syncNonFeatureBranchSteps(branchName string, repo *git.ProdRepo) (runstate.StepList, error) {
	hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branchName)
	if err != nil {
		return runstate.StepList{}, err
	}
	result := runstate.StepList{}
	if hasTrackingBranch {
		if repo.Config.PullBranchStrategy() == "rebase" {
			result.Append(&steps.RebaseBranchStep{BranchName: repo.Silent.TrackingBranchName(branchName)})
		} else {
			result.Append(&steps.MergeBranchStep{BranchName: repo.Silent.TrackingBranchName(branchName)})
		}
	}

	mainBranchName := repo.Config.MainBranch()
	hasUpstream, err := repo.Silent.HasRemote("upstream")
	if err != nil {
		return runstate.StepList{}, err
	}
	shouldSyncUpstream, err := repo.Config.ShouldSyncUpstream()
	if err != nil {
		return runstate.StepList{}, err
	}
	if mainBranchName == branchName && hasUpstream && shouldSyncUpstream {
		result.Append(&steps.FetchUpstreamStep{BranchName: mainBranchName})
		result.Append(&steps.RebaseBranchStep{BranchName: fmt.Sprintf("upstream/%s", mainBranchName)})
	}
	return result, nil
}

// provides the name of the branch that should be checked out after all sync steps run
func finalBranch(config syncConfig) string {
	if stringslice.Contains(config.branchesWithDeletedRemote, config.initialBranch) {
		return config.mainBranch
	} else {
		return config.initialBranch
	}
}
