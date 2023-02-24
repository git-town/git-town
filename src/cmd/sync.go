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
	"github.com/spf13/cobra"
)

type SyncConfig struct {
	branchesToSync   []string
	initialBranch    string
	syncBranchConfig SyncBranchConfig
}

type SyncBranchConfig struct {
	hasOrigin          bool
	hasUpstream        bool
	isOffline          bool
	mainBranch         string
	pushHook           bool
	shouldPushTags     bool
	shouldSyncUpstream bool
	syncStrategy       string
}

func syncCmd(repo *git.ProdRepo) *cobra.Command {
	// TODO: move to bottom of this function
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
			syncConfig, err := determineSyncConfig(allFlag, repo)
			if err != nil {
				cli.Exit(err)
			}
			stepList, err := syncSteps(syncConfig, repo)
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

func determineSyncConfig(allFlag bool, repo *git.ProdRepo) (*SyncConfig, error) {
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return nil, err
	}
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return nil, err
	}
	syncBranchConfig, err := determineSyncBranchConfig(repo, initialBranch)
	if err != nil {
		return nil, err
	}
	if hasOrigin && !isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return nil, err
		}
	}
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return nil, err
	}
	parentDialog := dialog.ParentBranches{}
	var branchesToSync []string
	var shouldPushTags bool
	if allFlag {
		branches, err := repo.Silent.LocalBranchesMainFirst()
		if err != nil {
			return nil, err
		}
		err = parentDialog.EnsureKnowsParentBranches(branches, repo)
		if err != nil {
			return nil, err
		}
		syncConfig.branchesToSync = branches
		syncBranchConfig.shouldPushTags = true
	} else {
		err = parentDialog.EnsureKnowsParentBranches([]string{syncConfig.initialBranch}, repo)
		if err != nil {
			return nil, err
		}
		syncConfig.branchesToSync = append(repo.Config.AncestorBranches(syncConfig.initialBranch), syncConfig.initialBranch)
	}
	return &syncConfig, nil
}

func determineSyncBranchConfig(repo *git.ProdRepo, initialBranch string) (SyncBranchConfig, error) {
	result := SyncBranchConfig{}
	var err error
	result.hasOrigin, err = repo.Silent.HasOrigin()
	if err != nil {
		return result, err
	}
	result.isOffline, err = repo.Config.IsOffline()
	if err != nil {
		return result, err
	}
	result.pushHook, err = repo.Config.PushHook()
	if err != nil {
		return result, err
	}
	result.mainBranch = repo.Config.MainBranch()
	result.hasUpstream, err = repo.Silent.HasRemote("upstream")
	if err != nil {
		return result, err
	}
	result.shouldSyncUpstream, err = repo.Config.ShouldSyncUpstream()
	if err != nil {
		return result, err
	}
	result.shouldPushTags = !repo.Config.IsFeatureBranch(initialBranch)
	return result, nil
}

// syncSteps provides the step list for the "git sync" command.
func syncSteps(config *SyncConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	for _, branch := range config.branchesToSync {
		steps, err := syncBranchSteps(branch, true, config.syncBranchConfig, repo)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	}
	result.Append(&steps.CheckoutBranchStep{Branch: config.initialBranch})
	if config.syncBranchConfig.hasOrigin && config.syncBranchConfig.shouldPushTags && !config.syncBranchConfig.isOffline {
		result.Append(&steps.PushTagsStep{})
	}
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	return result, err
}

// syncBranchSteps provides the steps to sync a particular branch.
func syncBranchSteps(branch string, pushBranch bool, config SyncBranchConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	isFeatureBranch := repo.Config.IsFeatureBranch(branch)
	if !config.hasOrigin && !isFeatureBranch {
		return runstate.StepList{}, nil
	}
	hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branch)
	if err != nil {
		return runstate.StepList{}, err
	}
	result.Append(&steps.CheckoutBranchStep{Branch: branch})
	if isFeatureBranch {
		steps, err := syncFeatureBranchSteps(branch, config, hasTrackingBranch, pushBranch, repo)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	} else {
		steps, err := syncNonFeatureBranchSteps(branch, config, hasTrackingBranch, pushBranch, repo)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	}
	return result, nil
}

// syncFeatureBranchSteps provides the steps to sync the feature branch with the given name.
func syncFeatureBranchSteps(branch string, config SyncBranchConfig, hasTrackingBranch bool, pushBranch bool, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	if hasTrackingBranch {
		steps, err := syncTrackingBranchSteps(repo.Silent.TrackingBranch(branch), config.syncStrategy)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	}
	parentSteps, err := syncParentSteps(repo.Config.ParentBranch(branch), config.syncStrategy)
	if err != nil {
		return runstate.StepList{}, err
	}
	result.AppendList(parentSteps)
	if pushBranch && config.hasOrigin && !config.isOffline {
		if !hasTrackingBranch {
			result.Append(&steps.CreateTrackingBranchStep{Branch: branch})
		} else {
			steps, err := pushFeatureBranchSteps(branch, config.syncStrategy, config.pushHook)
			if err != nil {
				return runstate.StepList{}, err
			}
			result.AppendList(steps)
		}
	}
	return result, nil
}

// syncNonFeatureBranchSteps provides the steps to sync the non-feature branch with the given name.
func syncNonFeatureBranchSteps(branch string, config SyncBranchConfig, hasTrackingBranch bool, pushBranch bool, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	if hasTrackingBranch {
		steps, err := syncTrackingBranchSteps(repo.Silent.TrackingBranch(branch), repo.Config.PullBranchStrategy())
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	}
	if config.mainBranch == branch && config.hasUpstream && config.shouldSyncUpstream {
		result.Append(&steps.FetchUpstreamStep{Branch: config.mainBranch})
		result.Append(&steps.RebaseBranchStep{Branch: fmt.Sprintf("upstream/%s", config.mainBranch)})
	}
	if pushBranch && config.hasOrigin && !config.isOffline {
		if !hasTrackingBranch {
			result.Append(&steps.CreateTrackingBranchStep{Branch: branch})
		} else {
			result.Append(&steps.PushBranchStep{Branch: branch})
		}
	}
	return result, nil
}

// syncTrackingBranchStep provides the steps to sync the given tracking branch into the current branch.
func syncTrackingBranchSteps(trackingBranch, syncStrategy string) (runstate.StepList, error) {
	switch syncStrategy {
	case "merge":
		return runstate.NewStepList(&steps.MergeBranchStep{Branch: trackingBranch}), nil
	case "rebase":
		return runstate.NewStepList(&steps.RebaseBranchStep{Branch: trackingBranch}), nil
	default:
		return runstate.StepList{}, fmt.Errorf("unknown syncStrategy value: %q", syncStrategy)
	}
}

// syncParentSteps provides the steps to sync the given parent branch into the current branch.
func syncParentSteps(parentBranch, syncStrategy string) (runstate.StepList, error) {
	switch syncStrategy {
	case "merge":
		return runstate.NewStepList(&steps.MergeBranchStep{Branch: parentBranch}), nil
	case "rebase":
		return runstate.NewStepList(&steps.RebaseBranchStep{Branch: parentBranch}), nil
	default:
		return runstate.StepList{}, fmt.Errorf("unknown syncStrategy value: %q", syncStrategy)
	}
}

func pushFeatureBranchSteps(branch, syncStrategy string, pushHook bool) (runstate.StepList, error) {
	switch syncStrategy {
	case "merge":
		return runstate.NewStepList(&steps.PushBranchStep{Branch: branch, NoPushHook: !pushHook}), nil
	case "rebase":
		return runstate.NewStepList(&steps.PushBranchStep{Branch: branch, ForceWithLease: true}), nil
	default:
		return runstate.StepList{}, fmt.Errorf("unknown syncStrategy value: %q", syncStrategy)
	}
}
