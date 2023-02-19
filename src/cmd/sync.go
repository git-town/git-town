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

type syncConfig struct {
	branchesToSync []string
	hasOrigin      bool
	initialBranch  string
	isOffline      bool
	shouldPushTags bool
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
			config, err := determineSyncConfig(allFlag, repo)
			if err != nil {
				cli.Exit(err)
			}
			stepList, err := syncSteps(config, repo)
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

func determineSyncConfig(allFlag bool, repo *git.ProdRepo) (*syncConfig, error) {
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return nil, err
	}
	isOffline, err := repo.Config.IsOffline()
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
		branchesToSync = branches
		shouldPushTags = true
	} else {
		err = parentDialog.EnsureKnowsParentBranches([]string{initialBranch}, repo)
		if err != nil {
			return nil, err
		}
		branchesToSync = append(repo.Config.AncestorBranches(initialBranch), initialBranch)
		shouldPushTags = !repo.Config.IsFeatureBranch(initialBranch)
	}
	return &syncConfig{
		branchesToSync: branchesToSync,
		hasOrigin:      hasOrigin,
		initialBranch:  initialBranch,
		isOffline:      isOffline,
		shouldPushTags: shouldPushTags,
	}, nil
}

// syncSteps provides the step list for the "git sync" command.
func syncSteps(config *syncConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	result := runstate.StepList{}
	for _, branch := range config.branchesToSync {
		steps, err := syncBranchSteps(branch, true, repo)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	}
	result.Append(&steps.CheckoutBranchStep{Branch: config.initialBranch})
	if config.hasOrigin && config.shouldPushTags && !config.isOffline {
		result.Append(&steps.PushTagsStep{})
	}
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, repo)
	return result, err
}

// syncBranchSteps provides the steps to sync a particular branch.
func syncBranchSteps(branch string, pushBranch bool, repo *git.ProdRepo) (runstate.StepList, error) {
	isFeatureBranch := repo.Config.IsFeatureBranch(branch)
	syncStrategy := repo.Config.SyncStrategy()
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return runstate.StepList{}, err
	}
	pushHook, err := repo.Config.PushHook()
	if err != nil {
		return runstate.StepList{}, err
	}
	result := runstate.StepList{}
	if !hasOrigin && !isFeatureBranch {
		return runstate.StepList{}, nil
	}
	result.Append(&steps.CheckoutBranchStep{Branch: branch})
	if isFeatureBranch {
		steps, err := syncFeatureBranchSteps(branch, repo)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	} else {
		steps, err := syncNonFeatureBranchSteps(branch, repo)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return runstate.StepList{}, err
	}
	if pushBranch && hasOrigin && !isOffline {
		hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branch)
		if err != nil {
			return runstate.StepList{}, err
		}
		if !hasTrackingBranch {
			result.Append(&steps.CreateTrackingBranchStep{Branch: branch})
			return result, nil
		}
		if !isFeatureBranch {
			result.Append(&steps.PushBranchStep{Branch: branch})
			return result, nil
		}
		steps, err := pushFeatureBranchSteps(branch, syncStrategy, pushHook)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	}
	return result, nil
}

func syncFeatureBranchSteps(branch string, repo *git.ProdRepo) (runstate.StepList, error) {
	syncStrategy := repo.Config.SyncStrategy()
	hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branch)
	if err != nil {
		return runstate.StepList{}, err
	}
	result := runstate.StepList{}
	if hasTrackingBranch {
		steps, err := syncTrackingBranchSteps(repo.Silent.TrackingBranch(branch), syncStrategy)
		if err != nil {
			return runstate.StepList{}, err
		}
		result.AppendList(steps)
	}
	steps, err := syncParentSteps(repo.Config.ParentBranch(branch), syncStrategy)
	if err != nil {
		return runstate.StepList{}, err
	}
	result.AppendList(steps)
	return result, nil
}

func syncNonFeatureBranchSteps(branch string, repo *git.ProdRepo) (runstate.StepList, error) {
	hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branch)
	if err != nil {
		return runstate.StepList{}, err
	}
	result := runstate.StepList{}
	if hasTrackingBranch {
		result, err = syncTrackingBranchSteps(repo.Silent.TrackingBranch(branch), repo.Config.PullBranchStrategy())
		if err != nil {
			return runstate.StepList{}, err
		}
	}
	mainBranch := repo.Config.MainBranch()
	hasUpstream, err := repo.Silent.HasRemote("upstream")
	if err != nil {
		return runstate.StepList{}, err
	}
	shouldSyncUpstream, err := repo.Config.ShouldSyncUpstream()
	if err != nil {
		return runstate.StepList{}, err
	}
	if mainBranch == branch && hasUpstream && shouldSyncUpstream {
		result.Append(&steps.FetchUpstreamStep{Branch: mainBranch})
		result.Append(&steps.RebaseBranchStep{Branch: fmt.Sprintf("upstream/%s", mainBranch)})
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
